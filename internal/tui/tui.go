package tui

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
)

type devicesMsg struct {
	devices []app.Device
	err     error
}
type preparedMsg struct {
	cmd            *exec.Cmd
	serial, preset string
	err            error
}
type endedMsg struct {
	err  error
	logs string
}
type Model struct {
	Filter    string
	Searching bool
	Search    textinput.Model
	Help      bool
	LogsView  bool
	Logs      string

	CommandView bool
	Offset      int

	A           *app.App
	Devices     []app.Device
	D, P, Focus int
	W, H        int
	Status      string
	SaveWarning string
	Busy        bool
	Editing     bool
	Delete      bool
	Field       int
	Inputs      []textinput.Model
	Extra       []string
	Creating    bool
}

func New(a *app.App, warning string) Model {
	m := Model{A: a, D: -1, W: 80, H: 24, Status: warning, Busy: true}
	for i, p := range a.Config.Presets {
		if p.Name == a.Config.LastPreset {
			m.P = i
		}
	}
	return m
}
func (m Model) discover() tea.Cmd {
	return func() tea.Msg { ds, err := m.A.Discover(context.Background()); return devicesMsg{ds, err} }
}
func (m Model) Init() tea.Cmd { return m.discover() }
func (m Model) serial() string {
	if m.D >= 0 && m.D < len(m.Devices) && !m.matches(m.Devices[m.D]) {
		return ""
	}
	if m.D >= 0 && m.D < len(m.Devices) && m.Devices[m.D].State == "device" {
		return m.Devices[m.D].Serial
	}
	return ""
}
func (m *Model) refresh(ds []app.Device) {
	m.Busy = false
	id := m.serial()
	if id == "" {
		id = m.A.Config.LastDevice
	}
	m.Devices = ds
	m.D = -1
	first := -1
	for i, d := range ds {
		if d.State == "device" && m.matches(d) {
			if first < 0 {
				first = i
			}
			if d.Serial == id {
				m.D = i
			}
		}
	}
	if m.D < 0 {
		m.D = first
	}
}
func (m *Model) move(delta int) {
	if m.Focus == 1 {
		n := len(m.A.Config.Presets)
		if n > 0 {
			m.P = (m.P + delta + n) % n
		}
		return
	}
	n := len(m.Devices)
	if n == 0 {
		return
	}
	i := m.D
	for k := 0; k < n; k++ {
		i = (i + delta + n) % n
		if m.Devices[i].State == "device" && m.matches(m.Devices[i]) {
			m.D = i
			return
		}
	}
}
func (m *Model) edit(create bool) {
	m.Editing = true
	m.Creating = create
	m.Field = 0
	m.Inputs = nil
	values := []string{"", "", ""}
	if !create {
		p := m.A.Config.Presets[m.P]
		parts := []string{}
		for _, s := range p.Args {
			parts = append(parts, app.Quote(s))
		}
		values = []string{p.Name, p.Description, strings.Join(parts, " ")}
	}
	for i, v := range values {
		t := textinput.New()
		t.SetValue(v)
		t.CharLimit = 4096
		t.Width = 60
		if i == 0 {
			t.Focus()
		}
		m.Inputs = append(m.Inputs, t)
	}
}
func (m *Model) persist() {
	m.SaveWarning = ""
	if err := m.A.Save(); err != nil {
		m.SaveWarning = "Enregistrement impossible : " + err.Error()
		m.Status = m.SaveWarning
	}
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.W = v.Width
		m.H = v.Height
	case devicesMsg:
		m.Busy = false
		if v.err != nil {
			m.Devices = nil
			m.D = -1
			m.Status = v.err.Error()
		} else {
			m.refresh(v.devices)
			if m.Status == "Rafraîchissement…" || m.Status == "" {
				m.Status = "Liste ADB actualisée"
			}
		}
	case preparedMsg:
		m.Busy = false
		if v.err != nil {
			m.Status = v.err.Error()
			return m, nil
		}
		m.A.Config.LastDevice = v.serial
		m.A.Config.LastPreset = v.preset
		m.persist()
		capture := &sessionCommand{cmd: v.cmd}
		return m, tea.Exec(capture, func(err error) tea.Msg { return endedMsg{err: err, logs: capture.tail.String()} })
	case endedMsg:
		m.Busy = true
		m.Logs = v.logs
		var interrupted *app.Interrupted
		if errors.As(v.err, &interrupted) && interrupted.Signal == syscall.SIGTERM {
			return m, tea.Quit
		}
		if v.err != nil {
			m.Status = "Échec scrcpy : " + v.err.Error() + " — l : diagnostic"
		} else {
			m.Status = "Session terminée ; Entrée pour relancer"
		}
		if m.SaveWarning != "" {
			m.Status = m.SaveWarning + " · " + m.Status
		}
		return m, m.discover()
	case tea.KeyMsg:
		key := v.String()
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if m.Busy {
			if key == "q" {
				return m, tea.Quit
			}
			return m, nil
		}
		if m.Delete {
			m.Delete = false
			if key == "o" {
				old := m.A.Config
				m.A.Config.Presets = append(append([]app.Preset{}, old.Presets[:m.P]...), old.Presets[m.P+1:]...)
				if err := m.A.Save(); err != nil {
					m.A.Config = old
					m.Status = err.Error()
				} else {
					m.P = 0
					m.Status = "Preset supprimé"
				}
			}
			return m, nil
		}
		if m.Editing {
			if key == "esc" {
				m.Editing = false
				return m, nil
			}
			if key == "tab" || key == "shift+tab" {
				m.Inputs[m.Field].Blur()
				delta := 1
				if key == "shift+tab" {
					delta = 2
				}
				m.Field = (m.Field + delta) % 3
				m.Inputs[m.Field].Focus()
				return m, textinput.Blink
			}
			if key == "enter" {
				args, err := app.ParseArgs(m.Inputs[2].Value())
				p := app.Preset{Name: strings.TrimSpace(m.Inputs[0].Value()), Description: m.Inputs[1].Value(), Args: args}
				if err == nil {
					err = app.ValidatePreset(p)
				}
				for i, other := range m.A.Config.Presets {
					if other.Name == p.Name && (m.Creating || i != m.P) {
						err = fmt.Errorf("ce nom existe déjà")
					}
				}
				if err != nil {
					m.Status = err.Error()
					return m, nil
				}
				old := m.A.Config
				ps := append([]app.Preset{}, old.Presets...)
				if m.Creating {
					ps = append(ps, p)
				} else {
					ps[m.P] = p
				}
				m.A.Config.Presets = ps
				if err = m.A.Save(); err != nil {
					m.A.Config = old
					m.Status = err.Error()
					return m, nil
				}
				if m.Creating {
					m.P = len(ps) - 1
				}
				m.Editing = false
				m.Status = "Preset enregistré"
				return m, nil
			}
			var cmd tea.Cmd
			m.Inputs[m.Field], cmd = m.Inputs[m.Field].Update(v)
			return m, cmd
		}

		if m.Help || m.LogsView {
			switch key {
			case "esc", "?", "l":
				m.Help = false
				m.LogsView = false
			case "down", "j":
				m.Offset++
			case "up", "k":
				m.Offset = max(0, m.Offset-1)
			case "q":
				return m, tea.Quit
			}
			return m, nil
		}
		if m.Searching {
			switch key {
			case "esc":
				m.Searching = false
				m.Filter = ""
				m.Search.SetValue("")
				m.reselect()
				return m, nil
			case "enter":
				m.Searching = false
				return m, nil
			}
			var cmd tea.Cmd
			m.Search, cmd = m.Search.Update(v)
			m.Filter = m.Search.Value()
			m.reselect()
			return m, cmd
		}
		if m.CommandView {
			switch key {
			case "esc", "c":
				m.CommandView = false
			case "down":
				m.Offset++
			case "up":
				m.Offset = max(0, m.Offset-1)
			case "q":
				return m, tea.Quit
			}
			return m, nil
		}
		if strings.HasPrefix(key, "/") && len(key) > 1 {
			m.Searching = true
			m.Search = textinput.New()
			m.Filter = strings.TrimPrefix(key, "/")
			m.Search.SetValue(m.Filter)
			m.Search.Focus()
			m.reselect()
			return m, textinput.Blink
		}
		switch key {
		case "/":
			m.Searching = true
			m.Search = textinput.New()
			m.Search.SetValue(m.Filter)
			m.Search.Focus()
			return m, textinput.Blink
		case "esc":
			m.Filter = ""
			m.reselect()
		case "?":
			m.Help = true
			m.Offset = 0
		case "l":
			m.LogsView = true
			m.Offset = 0
		case "home":
			m.jump(false)
		case "end":
			m.jump(true)
		case "pgup":
			m.page(-1)
		case "pgdown":
			m.page(1)
		case "c":
			m.CommandView = true
			m.Offset = 0
		case "q":
			return m, tea.Quit
		case "tab", "left", "right":
			m.Focus = 1 - m.Focus
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "r":
			if !m.Busy {
				m.Busy = true
				m.Status = "Rafraîchissement…"
				return m, m.discover()
			}
		case "n":
			m.edit(true)
			return m, textinput.Blink
		case "e":
			if len(m.A.Config.Presets) > 0 {
				m.edit(false)
				return m, textinput.Blink
			}
		case "d":
			if len(m.A.Config.Presets) > 0 {
				m.Delete = true
			}
		case "enter":
			if m.Busy {
				return m, nil
			}
			if m.serial() == "" || len(m.A.Config.Presets) == 0 {
				m.Status = "Sélectionnez un appareil disponible et un preset"
				return m, nil
			}
			serial := m.serial()
			p := m.A.Config.Presets[m.P]
			m.Busy = true
			m.Status = "Vérification de l'appareil…"
			return m, func() tea.Msg {
				cmd, err := m.A.Prepare(context.Background(), serial, p, m.Extra)
				return preparedMsg{cmd, serial, p.Name, err}
			}
		}
	}
	return m, nil
}

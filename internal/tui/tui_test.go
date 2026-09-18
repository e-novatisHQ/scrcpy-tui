package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
)

func key(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }
func TestNavigationRefreshAndCancel(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	m.refresh([]app.Device{{Serial: "A", State: "device"}, {Serial: "B", State: "unauthorized"}, {Serial: "C", State: "device"}})
	m.move(1)
	if m.serial() != "C" {
		t.Fatal(m.serial())
	}
	m.refresh([]app.Device{{Serial: "C", State: "device"}, {Serial: "A", State: "device"}})
	if m.serial() != "C" {
		t.Fatal("selection lost")
	}
	v, _ := m.Update(key("n"))
	m = v.(Model)
	if !m.Editing {
		t.Fatal("editor")
	}
	v, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = v.(Model)
	if m.Editing || len(a.Config.Presets) != 3 {
		t.Fatal("cancel")
	}
	v, _ = m.Update(key("d"))
	m = v.(Model)
	v, _ = m.Update(key("n"))
	m = v.(Model)
	if len(a.Config.Presets) != 3 {
		t.Fatal("delete refusal")
	}
}
func TestDimensions(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	for i := 0; i < 40; i++ {
		m.Devices = append(m.Devices, app.Device{Serial: strings.Repeat("界", 30), State: "device"})
	}
	m.D = 39
	for _, w := range []int{30, 60, 89, 90, 120} {
		for _, h := range []int{12, 18, 24} {
			m.W = w
			m.H = h
			view := m.View()
			lines := strings.Split(view, "\n")
			if len(lines) > h {
				t.Fatal(w, h, len(lines))
			}
			for _, l := range lines {
				if ansi.StringWidth(l) > w {
					t.Fatal(w, l)
				}
			}
			if !strings.Contains(view, "Entrée") {
				t.Fatal("footer hidden", w, h)
			}
		}
	}
}

func TestSearchPagingAndFocus(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	for i := 0; i < 42; i++ {
		m.Devices = append(m.Devices, app.Device{Serial: fmt.Sprintf("10.0.0.1:%d", 6201+i), Model: "Dongle G 4K", State: "device"})
	}
	m.D = 0
	m.Busy = false
	m.Filter = "6242"
	m.reselect()
	if m.serial() != "10.0.0.1:6242" {
		t.Fatal(m.serial())
	}
	m.Filter = "absent"
	m.reselect()
	if m.serial() != "" {
		t.Fatal("hidden selection")
	}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("empty filter launches")
	}
	m.Filter = ""
	m.reselect()
	m.jump(true)
	if m.D != 41 {
		t.Fatal(m.D)
	}
	m.jump(false)
	m.page(1)
	if m.D <= 0 {
		t.Fatal("page no movement")
	}
	for _, wh := range [][2]int{{30, 12}, {60, 12}, {89, 18}, {90, 18}, {120, 48}} {
		m.W = wh[0]
		m.H = wh[1]
		v := ansi.Strip(m.View())
		if !strings.Contains(v, "> ") || !strings.Contains(v, "Presets") {
			t.Fatal("selection missing", wh, v)
		}
		if strings.Contains(v, "[1/42]") {
			t.Fatal("counter on row")
		}
	}
	v, _ := m.Update(key("?"))
	m = v.(Model)
	if !m.Help {
		t.Fatal("help")
	}
	v, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = v.(Model)
	if m.Help {
		t.Fatal("help cancel")
	}
}
func TestBoundedCapture(t *testing.T) {
	b := &tailBuffer{}
	b.Write([]byte(strings.Repeat("x", 40000)))
	b.Write([]byte("ERROR invalid option"))
	s := b.String()
	if len(s) > 32768 || !strings.HasSuffix(s, "ERROR invalid option") {
		t.Fatal(len(s))
	}
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	v, _ := m.Update(endedMsg{err: fmt.Errorf("exit status 7"), logs: s})
	m = v.(Model)
	if m.Logs != s || !strings.Contains(m.Status, "l : diagnostic") {
		t.Fatal("missing diagnostic")
	}
}

func TestGroupedSearchInput(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	m.refresh([]app.Device{{Serial: "host:6201", State: "device"}, {Serial: "host:6242", State: "device"}})
	v, _ := m.Update(key("/6242"))
	m = v.(Model)
	if !m.Searching || m.serial() != "host:6242" {
		t.Fatal(m.Filter, m.serial())
	}
	v, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = v.(Model)
	if m.Searching || m.Filter != "" {
		t.Fatal("filter not cancelled")
	}
}
func TestEmptyAndAnomalyRendering(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	state, _ := m.Update(devicesMsg{err: fmt.Errorf("ADB unavailable")})
	m = state.(Model)
	v := ansi.Strip(m.View())
	if !strings.Contains(v, "Lancement bloqué") || !strings.Contains(v, "ADB unavailable") {
		t.Fatal(v)
	}
	m.refresh([]app.Device{{Serial: "USB", State: "unauthorized", Model: "TV"}})
	v = ansi.Strip(m.View())
	if !strings.Contains(v, "unauthorized") || m.serial() != "" {
		t.Fatal(v)
	}
	m.Logs = "\x1b[2JERROR option inconnue\nNext line"
	m.LogsView = true
	v = m.View()
	if strings.Contains(v, "\x1b[2J") || !strings.Contains(v, "ERROR option inconnue") {
		t.Fatal("unsafe diagnostic")
	}
}

func TestDiscoveryLocksLaunch(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	if !m.Busy {
		t.Fatal("initial discovery not locked")
	}
	v, _ := m.Update(devicesMsg{devices: []app.Device{{Serial: "USB", State: "device"}}})
	m = v.(Model)
	if m.Busy {
		t.Fatal("discovery never unlocks")
	}
	v, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = v.(Model)
	if !m.Busy || cmd == nil {
		t.Fatal("prepare not locked")
	}
	v, _ = m.Update(key("n"))
	if v.(Model).Editing {
		t.Fatal("edit during prepare")
	}
	v, _ = m.Update(endedMsg{})
	if !v.(Model).Busy {
		t.Fatal("return discovery not locked")
	}
}
func TestPageDoesNotWrap(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "")
	for i := 0; i < 12; i++ {
		m.Devices = append(m.Devices, app.Device{Serial: fmt.Sprint(i), State: "device"})
	}
	m.D = 0
	m.W = 100
	m.page(1)
	if m.D != 11 {
		t.Fatal("page did not reach end", m.D)
	}
	m.page(1)
	if m.D != 11 {
		t.Fatal("page wrapped")
	}
	m.page(-1)
	if m.D != 0 {
		t.Fatal("page did not reach start")
	}
	m.Focus = 1
	m.P = 0
	m.page(1)
	if m.P != 2 {
		t.Fatal("preset page wraps", m.P)
	}
}
func TestErrorVisibleInSmallTerminal(t *testing.T) {
	a, _ := app.New(filepath.Join(t.TempDir(), "p.json"))
	m := New(a, "ADB : commande introuvable")
	m.W = 40
	m.H = 12
	if !strings.Contains(ansi.Strip(m.View()), "ADB : commande introuvable") {
		t.Fatal("error hidden")
	}
}

func TestSaveFailureRemainsAfterSession(t *testing.T) {
	a, _ := app.New(t.TempDir())
	m := New(a, "")
	m.persist()
	if m.SaveWarning == "" {
		t.Fatal("write failure ignored")
	}
	v, _ := m.Update(endedMsg{})
	m = v.(Model)
	if !strings.Contains(m.Status, "Enregistrement impossible") {
		t.Fatal("write failure erased at session end")
	}
}

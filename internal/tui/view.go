package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
)

var accent = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
var warning = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
var muted = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
var selected = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
var focused = lipgloss.NewStyle().Foreground(lipgloss.Color("232")).Background(lipgloss.Color("81")).Bold(true)

func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return ' '
		}
		return r
	}, s)
}
func line(s string, w int) string { return ansi.Truncate(clean(s), max(1, w), "…") }
func (m Model) matches(d app.Device) bool {
	return strings.Contains(strings.ToLower(d.Serial+" "+d.Model), strings.ToLower(m.Filter))
}
func (m Model) visible() []int {
	var ids []int
	for i, d := range m.Devices {
		if m.matches(d) {
			ids = append(ids, i)
		}
	}
	return ids
}
func (m *Model) reselect() {
	if m.serial() != "" {
		return
	}
	m.D = -1
	for _, i := range m.visible() {
		if m.Devices[i].State == "device" {
			m.D = i
			return
		}
	}
}
func (m *Model) jump(end bool) {
	if m.Focus == 1 {
		m.P = 0
		if end {
			m.P = max(0, len(m.A.Config.Presets)-1)
		}
		return
	}
	ids := m.visible()
	if end {
		for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
			ids[i], ids[j] = ids[j], ids[i]
		}
	}
	for _, i := range ids {
		if m.Devices[i].State == "device" {
			m.D = i
			return
		}
	}
}
func (m *Model) page(direction int) {
	step := max(1, m.H-12)
	if m.W < 90 {
		step = max(1, (m.H-12)/2)
	}
	if m.Focus == 1 {
		m.P = max(0, min(len(m.A.Config.Presets)-1, m.P+direction*step))
		return
	}
	var ids []int
	pos := 0
	for _, i := range m.visible() {
		if m.Devices[i].State == "device" {
			if i == m.D {
				pos = len(ids)
			}
			ids = append(ids, i)
		}
	}
	if len(ids) > 0 {
		m.D = ids[max(0, min(len(ids)-1, pos+direction*step))]
	}
}

func span(n, sel, budget int) (int, int) {
	start := max(0, sel-budget+1)
	return start, min(n, start+budget)
}
func (m Model) command() string {
	if m.serial() == "" || len(m.A.Config.Presets) == 0 {
		return "Sélection incomplète"
	}
	args, err := m.A.Plan(m.serial(), m.A.Config.Presets[m.P], m.Extra)
	if err != nil {
		return err.Error()
	}
	return app.Command(args)
}
func (m Model) ready() string {
	if m.Busy {
		return "Vérification en cours…"
	}
	if m.serial() == "" {
		return "Lancement bloqué : aucun appareil disponible"
	}
	if len(m.A.Config.Presets) == 0 {
		return "Lancement bloqué : créez un preset (n)"
	}
	return "Prêt à lancer · Entrée"
}
func (m Model) selection(s string, on, active bool, w int) string {
	s = line(s, w)
	if on {
		if active {
			return focused.Width(w).Render(s)
		}
		return selected.Render(s)
	}
	return muted.Render(s)
}
func (m Model) detail(title, text string, w int) string {
	parts := strings.Split(ansi.Strip(text), "\n")
	for i, s := range parts {
		parts[i] = clean(s)
	}
	ls := strings.Split(ansi.Hardwrap(strings.Join(parts, "\n"), w, true), "\n")
	budget := m.H - 3
	start := min(m.Offset, max(0, len(ls)-budget))
	end := min(len(ls), start+budget)
	out := []string{accent.Render(line(title, w))}
	out = append(out, ls[start:end]...)
	out = append(out, line("↑/↓ défiler · Échap retour", w))
	return strings.Join(out, "\n")
}
func (m Model) View() string {
	if m.W < 30 || m.H < 12 {
		return "Terminal trop petit (30×12).\nAgrandir · q quitter"
	}
	w := m.W - 2
	if m.Help {
		return m.detail("Aide clavier", "NAVIGATION\nTab / ← / → : changer de liste\n↑ / ↓ ou k / j : sélectionner\nHome / End : premier / dernier\nPage↑ / Page↓ : avancer par page\n/ : rechercher adresse, port ou modèle\nEntrée : terminer la recherche\nÉchap : effacer le filtre\n\nSESSION\nEntrée : lancer appareil + preset sélectionnés\nr : rafraîchir ADB\nc : commande complète\nl : dernières sorties scrcpy\nq / Ctrl+C : quitter\n\nPRESETS (quel que soit le focus)\nn : créer un preset\ne : modifier le preset sélectionné\nd : supprimer le preset sélectionné\no : confirmer la suppression\nToute autre touche : annuler la suppression\n\nÉDITEUR\nTab / Shift+Tab : changer de champ\nEntrée : enregistrer\nÉchap : annuler", w)
	}
	if m.LogsView {
		text := m.Logs
		if text == "" {
			text = "Aucune sortie scrcpy pour cette session."
		}
		return m.detail("Diagnostic scrcpy — dernières sorties", text, w)
	}
	if m.CommandView {
		return m.detail("Commande exacte", m.command(), w)
	}
	rows := []string{accent.Render("scrcpy-tui · " + line(m.ready(), w-13))}
	if m.Editing {
		title := "Modifier le preset"
		if m.Creating {
			title = "Créer un preset"
		}
		rows = append(rows, title, line(m.Status, w))
		for i, t := range m.Inputs {
			t.Width = w - 3
			rows = append(rows, []string{"Nom", "Description", "Arguments scrcpy"}[i], t.View())
		}
		rows = append(rows, "Tab champ · Entrée sauver", "Échap annuler")
		return fit(rows, w, m.H)
	}
	name := "aucun preset"
	if len(m.A.Config.Presets) > 0 {
		name = m.A.Config.Presets[m.P].Name
	}
	rows = append(rows, line("Cible : "+m.serial()+" · "+name, w), muted.Render(line(m.Status, w)))
	filter := "/ rechercher · ? aide"
	if m.Filter != "" {
		filter = "Filtre : " + m.Filter + " · Échap effacer"
	}
	if m.Searching {
		m.Search.Width = w - 14
		filter = "Recherche : " + m.Search.View()
	}

	rows = append(rows, filter)
	ids := m.visible()
	sel := -1
	for j, i := range ids {
		if i == m.D {
			sel = j
		}
	}
	wide := m.W >= 90
	budget := max(1, m.H-12)
	if !wide {
		budget = max(1, (m.H-12)/2)
	}
	start, end := span(len(ids), sel, budget)
	dt := "Appareils"
	pt := "Presets"
	if m.Focus == 0 {
		dt += " [focus]"
	} else {
		pt += " [focus]"
	}
	left := []string{accent.Render(dt)}
	left = append(left, muted.Render(fmt.Sprintf("Sélection %d/%d · visibles %d–%d/%d", max(0, sel+1), len(ids), min(len(ids), start+1), end, len(m.Devices))))
	lw := w
	if wide {
		lw = (w - 3) * 2 / 3
	}
	if len(ids) == 0 {
		left = append(left, "Aucun appareil correspondant")
	} else {
		for _, i := range ids[start:end] {
			d := m.Devices[i]
			mark := "  "
			if i == m.D {
				mark = "> "
			}
			state := ""
			if d.State != "device" {
				state = " [" + d.State + "]"
			}
			label := mark + d.Serial + state
			if lw >= 50 {
				label = fmt.Sprintf("%s%-24s %s %s", mark, d.Serial, state, d.Model)
			}
			rendered := m.selection(label, i == m.D, m.Focus == 0, lw)
			if d.State != "device" {
				rendered = warning.Render(line(label, lw))
			}
			left = append(left, rendered)
		}
	}
	rw := w
	if wide {
		rw = w - lw - 3
	}
	right := []string{accent.Render(pt)}
	ps, pe := span(len(m.A.Config.Presets), m.P, budget)
	for i := ps; i < pe; i++ {
		p := m.A.Config.Presets[i]
		mark := "  "
		if i == m.P {
			mark = "> "
		}
		right = append(right, m.selection(mark+p.Name, i == m.P, m.Focus == 1, rw))
		if i == m.P && (wide || m.H >= 18) {
			for _, s := range descriptionLines(p.Description, rw) {
				right = append(right, muted.Render(s))
			}
		}
	}
	if len(m.A.Config.Presets) == 0 {
		right = append(right, "Aucun preset · n créer")

	}
	if wide {
		for i := 0; i < max(len(left), len(right)); i++ {
			a, b := "", ""
			if i < len(left) {
				a = left[i]
			}
			if i < len(right) {
				b = right[i]
			}
			rows = append(rows, lipgloss.NewStyle().Width(lw).Render(ansi.Truncate(a, lw, "…"))+" │ "+b)
		}
	} else {
		rows = append(rows, left...)
		rows = append(rows, right...)
	}
	foot := []string{line("Commande : "+m.command(), w), "Tab liste · Entrée lancer", "r actualiser · / chercher · c commande · l logs · ? aide · q quitter", "Presets : n créer · e modifier · d supprimer"}
	if w < 65 {
		foot = []string{line("Commande : "+m.command(), w), "Tab liste · Entrée lancer", "r / c l · ? aide · q quitter", "Presets : n/e/d · ? aide"}
	}
	if m.H < 18 {
		foot = []string{line("Commande : "+m.command(), w), "Tab liste · Entrée lancer", "r / c l n e d · ? · q quitter"}
	}
	if m.Delete {
		foot = []string{line("Supprimer « "+name+" » ?", w), "o supprimer · autre annuler"}
	}
	available := m.H - len(foot)
	if len(rows) > available {
		rows = rows[:available]
	}
	rows = append(rows, foot...)
	return fit(rows, w, m.H)
}
func fit(rows []string, w, h int) string {
	if len(rows) > h {
		rows = rows[:h]
	}
	for i, s := range rows {
		rows[i] = ansi.Truncate(s, w, "…")
	}
	return strings.Join(rows, "\n")
}

func descriptionLines(text string, w int) []string {
	var rows []string
	current := ""
	for _, part := range strings.Split(clean(text), " · ") {
		candidate := part
		if current != "" {
			candidate = current + " · " + part
		}
		if ansi.StringWidth(candidate) > w && current != "" {
			rows = append(rows, current)
			current = part
		} else {
			current = candidate
		}
	}
	rows = append(rows, current)
	var out []string
	for _, r := range rows {
		out = append(out, strings.Split(ansi.Hardwrap(r, w, true), "\n")...)
	}
	return out
}

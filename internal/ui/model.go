package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	projects  []string
	selected  int
	width     int
	height    int
	searching bool
	notice    string
}

func New() tea.Model {
	return newModel()
}

func newModel() model {
	return model{
		projects: []string{
			"Forge",
			"phub",
			"ScreenBot",
			"Website",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		return m.updateKey(msg.String())
	}

	return m, nil
}

func (m model) updateKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case "down", "j":
		if m.selected < len(m.projects)-1 {
			m.selected++
		}
	case "/":
		m.searching = true
		m.notice = "Search is coming in the next milestone."
	case "esc", "escape":
		m.searching = false
		m.notice = ""
	case "enter":
		m.notice = fmt.Sprintf("Opening %s is planned for a later milestone.", m.projects[m.selected])
	}

	return m, nil
}

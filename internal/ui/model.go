package ui

import (
	"context"
	"fmt"
	"slices"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	projects         []Project
	shell            string
	selected         int
	width            int
	height           int
	searching        bool
	notice           string
	theme            theme
	themeMenu        bool
	themeSelection   int
	loadProjects     projectLoader
	scope            ProjectScope
	scopeSelection   int
	startup          bool
	refreshing       bool
	terminal         terminalSession
	terminalProject  string
	startingTerminal bool
	terminalContext  context.Context
	startTerminal    terminalStarter
}

type projectsRefreshedMsg struct {
	projects []Project
	err      error
}

func New(projects []Project, shell string) tea.Model {
	return newModel(projects, shell)
}

func NewWithRefresh(projects []Project, shell string, refresh func() ([]Project, error)) tea.Model {
	model := newModel(projects, shell)
	model.loadProjects = func(ProjectScope) ([]Project, error) {
		return refresh()
	}
	return model
}

func newModel(projects []Project, shell string) model {
	return model{
		projects:        slices.Clone(projects),
		shell:           shell,
		theme:           defaultTheme,
		terminalContext: context.Background(),
		startTerminal:   defaultTerminalStarter,
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
		if m.terminal != nil {
			return m, resizeTerminal(m.terminal, msg.Width, msg.Height)
		}
	case tea.KeyPressMsg:
		return m.updateKey(msg)
	case terminalStartedMsg:
		m.startingTerminal = false
		if msg.err != nil {
			m.notice = fmt.Sprintf("Could not open terminal: %v", msg.err)
			return m, nil
		}
		m.terminal = msg.session
		m.terminalProject = msg.projectName
		m.notice = ""
		return m, readTerminal(msg.session)
	case terminalOutputMsg:
		return m.handleTerminalOutput(msg)
	case terminalResizeMsg:
		if msg.session == m.terminal && msg.err != nil {
			m.notice = fmt.Sprintf("Could not resize terminal: %v", msg.err)
		}
	case projectsRefreshedMsg:
		m.refreshing = false
		if msg.err != nil {
			if m.startup {
				m.notice = fmt.Sprintf("Could not load projects: %v", msg.err)
			} else {
				m.notice = fmt.Sprintf("Could not refresh projects: %v", msg.err)
			}
			return m, nil
		}
		m.projects = slices.Clone(msg.projects)
		if len(m.projects) == 0 {
			m.selected = 0
		} else if m.selected >= len(m.projects) {
			m.selected = len(m.projects) - 1
		}
		if m.startup {
			m.startup = false
			m.notice = fmt.Sprintf("Loaded %d projects.", len(m.projects))
		} else {
			m.notice = fmt.Sprintf("Refreshed %d projects.", len(m.projects))
		}
	}

	return m, nil
}

func (m model) updateKey(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.terminal != nil {
		if err := m.terminal.SendKey(key.Key()); err != nil {
			m.notice = fmt.Sprintf("Could not send terminal input: %v", err)
		}
		return m, nil
	}
	if m.startingTerminal {
		return m, nil
	}
	if m.startup {
		return m.updateScopeMenu(key.Keystroke())
	}
	if m.themeMenu {
		return m.updateThemeMenu(key.Keystroke())
	}

	switch key.Keystroke() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "ctrl+p":
		m.themeMenu = true
		m.notice = ""
		return m, nil
	case "R", "r", "shift+r":
		return m.beginProjectLoad(m.scope)
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
		return m.beginTerminal()
	}

	return m, nil
}

func (m model) updateThemeMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.themeSelection > 0 {
			m.themeSelection--
		}
		return m, nil
	case "down", "j":
		if m.themeSelection < len(themes)-1 {
			m.themeSelection++
		}
		return m, nil
	case "enter":
		m.theme = themes[m.themeSelection]
		m.themeMenu = false
		m.notice = fmt.Sprintf("Theme: %s", m.theme.name)
		return m, nil
	case "esc", "escape":
		m.themeMenu = false
		m.notice = ""
		return m, nil
	}

	return m, nil
}

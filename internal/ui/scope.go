package ui

import tea "charm.land/bubbletea/v2"

type ProjectScope int

const (
	ScopeGitHubOnly ProjectScope = iota
	ScopeAllLocal
)

type projectLoader func(ProjectScope) ([]Project, error)

var projectScopeOptions = [...]struct {
	scope ProjectScope
	label string
}{
	{scope: ScopeGitHubOnly, label: "GitHub only (default)"},
	{scope: ScopeAllLocal, label: "All local projects"},
}

func NewWithLoader(shell string, load func(ProjectScope) ([]Project, error)) tea.Model {
	return newStartupModel(shell, load)
}

func newStartupModel(shell string, load projectLoader) model {
	model := newModel(nil, shell)
	model.loadProjects = load
	model.startup = true
	return model
}

func (m model) beginProjectLoad(scope ProjectScope) (tea.Model, tea.Cmd) {
	if m.loadProjects == nil {
		m.notice = "Project loading is unavailable."
		return m, nil
	}
	if m.refreshing {
		return m, nil
	}

	m.scope = scope
	m.refreshing = true
	if m.startup {
		m.notice = "Loading projects..."
	} else {
		m.notice = "Refreshing projects..."
	}
	load := m.loadProjects
	return m, func() tea.Msg {
		projects, err := load(scope)
		return projectsRefreshedMsg{projects: projects, err: err}
	}
}

func (m model) updateScopeMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	if m.refreshing {
		return m, nil
	}

	switch key {
	case "up", "k":
		if m.scopeSelection > 0 {
			m.scopeSelection--
		}
	case "down", "j":
		if m.scopeSelection < len(projectScopeOptions)-1 {
			m.scopeSelection++
		}
	case "enter":
		return m.beginProjectLoad(projectScopeOptions[m.scopeSelection].scope)
	}

	return m, nil
}

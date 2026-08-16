package ui

import (
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
)

const compactWidth = 64

var Version = "dev"

func (m model) View() tea.View {
	view := tea.NewView(m.paintBackground(m.render()))
	view.AltScreen = true
	view.BackgroundColor = m.theme.background
	switch m.theme.mode {
	case foregroundMode, comboMode:
		view.ForegroundColor = m.theme.foreground
	}
	if m.terminal != nil && !m.themeMenu {
		x, y := m.terminal.CursorPosition()
		view.Cursor = tea.NewCursor(x, y+1)
	}
	return view
}

func (m model) render() string {
	if m.terminal != nil && !m.themeMenu {
		return m.terminalView()
	}
	if m.startup {
		return m.startupView()
	}
	if m.width > 0 && m.width < compactWidth {
		return m.compactView()
	}

	return m.standardView()
}

func (m model) startupView() string {
	var content strings.Builder

	content.WriteString("phub\n\n")
	if m.width > 0 && m.width < 36 {
		content.WriteString("Load project scope:\n\n")
	} else {
		content.WriteString("Choose which projects to load:\n\n")
	}
	for index, option := range projectScopeOptions {
		if index == m.scopeSelection {
			content.WriteString("> ")
		} else {
			content.WriteString("  ")
		}
		content.WriteString(option.label)
		content.WriteByte('\n')
	}
	if m.notice != "" {
		content.WriteByte('\n')
		content.WriteString(m.notice)
		content.WriteByte('\n')
	}
	content.WriteString("\nup/down | enter load | q quit · phub ")
	content.WriteString(Version)
	content.WriteByte('\n')

	return content.String()
}

func (m model) standardView() string {
	var content strings.Builder

	content.WriteString("phub\n\n")
	if m.searching {
		content.WriteString("Search: coming in a future milestone (Esc to return)\n\n")
	} else {
		content.WriteString("Search projects... (press /)\n\n")
	}
	if m.themeMenu {
		m.writeThemeMenu(&content)
	} else {
		m.writeProjects(&content)
	}
	if m.notice != "" {
		content.WriteByte('\n')
		content.WriteString(m.notice)
		content.WriteByte('\n')
	}
	content.WriteString("\nj/k arrows | enter open | ctrl+p themes | R refresh | q quit · phub ")
	content.WriteString(Version)
	content.WriteByte('\n')

	return content.String()
}

func (m model) writeThemeMenu(content *strings.Builder) {
	switch {
	case m.width > 0 && m.width < 36:
		content.WriteString("Theme: arrows + Enter / Esc\n\n")
	case m.width > 0 && m.width < compactWidth:
		content.WriteString("Theme: arrows, Enter select, Esc cancel\n\n")
	default:
		content.WriteString("Choose a theme (up/down, Enter select, Esc cancel)\n\n")
	}

	start, end := m.themeRange()
	if start > 0 {
		content.WriteString("  ^ more\n")
	}
	for index := start; index < end; index++ {
		if index == m.themeSelection {
			content.WriteString("> ")
		} else {
			content.WriteString("  ")
		}
		content.WriteString(themes[index].name)
		content.WriteByte('\n')
	}
	if end < len(themes) {
		content.WriteString("  v more\n")
	}
}

func (m model) themeRange() (int, int) {
	visible := len(themes)
	if m.height > 0 {
		visible = m.height - 13
		if visible < 1 {
			visible = 1
		}
		if visible > len(themes) {
			visible = len(themes)
		}
	}

	start := m.themeSelection - visible/2
	if start < 0 {
		start = 0
	}
	if start+visible > len(themes) {
		start = len(themes) - visible
	}

	return start, start + visible
}

func (m model) compactView() string {
	var content strings.Builder

	content.WriteString("phub\n\n")
	if m.searching {
		content.WriteString("Search planned (Esc)\n\n")
	}
	if m.themeMenu {
		m.writeThemeMenu(&content)
	} else {
		m.writeProjects(&content)
	}
	if m.notice != "" && (m.width == 0 || utf8.RuneCountInString(m.notice) <= m.width) {
		content.WriteByte('\n')
		content.WriteString(m.notice)
		content.WriteByte('\n')
	}
	if m.themeMenu {
		content.WriteString("\nup/down | enter select | esc cancel\n")
	} else if m.width > 0 && m.width < 36 {
		content.WriteString("\nj/k | enter | R | ctrl+p | q\n")
	} else {
		content.WriteString("\nj/k | enter | R refresh | ctrl+p | q\n")
	}

	return content.String()
}

func (m model) writeProjects(content *strings.Builder) {
	if len(m.projects) == 0 {
		content.WriteString("No projects discovered.\n")
		return
	}

	names := projectDisplayNames(m.projects)

	visible := len(m.projects)
	if m.height > 0 {
		available := m.height - 7
		if available < 1 {
			available = 1
		}
		if available < visible {
			visible = available
		}
	}

	start := 0
	if m.selected >= visible {
		start = m.selected - visible + 1
	}

	end := start + visible
	if end > len(m.projects) {
		end = len(m.projects)
		start = end - visible
	}

	for index := start; index < end; index++ {
		if index == m.selected {
			content.WriteString("> ")
		} else {
			content.WriteString("  ")
		}
		content.WriteString(names[index])
		content.WriteByte('\n')
	}

	if start > 0 || end < len(m.projects) {
		content.WriteString("  ...\n")
	}
}

package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

const compactWidth = 64

func (m model) View() tea.View {
	view := tea.NewView(m.render())
	view.AltScreen = true
	return view
}

func (m model) render() string {
	if m.width > 0 && m.width < compactWidth {
		return m.compactView()
	}

	return m.standardView()
}

func (m model) standardView() string {
	var content strings.Builder

	content.WriteString("phub\n\n")
	if m.searching {
		content.WriteString("Search: coming in a future milestone (Esc to return)\n\n")
	} else {
		content.WriteString("Search projects... (press /)\n\n")
	}

	m.writeProjects(&content)
	if m.notice != "" {
		content.WriteByte('\n')
		content.WriteString(m.notice)
		content.WriteByte('\n')
	}
	content.WriteString("\nj/k or arrows: navigate | enter: open | /: search | q: quit\n")

	return content.String()
}

func (m model) compactView() string {
	var content strings.Builder

	content.WriteString("phub\n\n")
	if m.searching {
		content.WriteString("Search planned (Esc)\n\n")
	}

	m.writeProjects(&content)
	if m.notice == "" {
		if m.width > 0 && m.width < 36 {
			content.WriteString("\nj/k | enter | q\n")
		} else {
			content.WriteString("\nj/k navigate | enter open | q quit\n")
		}

		return content.String()
	}

	content.WriteByte('\n')
	if m.searching {
		content.WriteString("Search planned\n")
	} else {
		content.WriteString("Open planned\n")
	}
	if m.width > 0 && m.width < 36 {
		content.WriteString("\nj/k | enter | q\n")
	} else {
		content.WriteString("\nj/k navigate | enter open | q quit\n")
	}

	return content.String()
}

func (m model) writeProjects(content *strings.Builder) {
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
		content.WriteString(m.projects[index])
		content.WriteByte('\n')
	}

	if start > 0 || end < len(m.projects) {
		content.WriteString("  ...\n")
	}
}

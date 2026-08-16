package ui

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
)

const terminalChromeRows = 2

const (
	terminalBold = "\x1b[1m"
	terminalDim  = "\x1b[2m"
	terminalTone = "\x1b[22m"
)

func terminalViewportDimensions(width int, height int) (int, int) {
	width, height = terminalDimensions(width, height)
	if height <= terminalChromeRows {
		height = 1
	} else {
		height -= terminalChromeRows
	}
	return width, height
}

func (m model) terminalView() string {
	width, _ := terminalDimensions(m.width, m.height)
	header := terminalChromeLine(width, "╭─", terminalHeaderLabel(m.terminalProject, terminalShellName(m.shell), width-3), "╮", terminalBold)
	footer := terminalChromeLine(width, "╰─", " Ctrl+D return · Ctrl+C interrupt ", "╯", terminalDim)
	return strings.Join([]string{header, m.terminal.Render(), footer}, "\n")
}

func terminalChromeLine(width int, prefix string, label string, suffix string, style string) string {
	if width < 1 {
		return ""
	}
	available := width - ansi.StringWidth(prefix) - ansi.StringWidth(suffix)
	label = fitTerminalLabel(label, available)
	fill := width - ansi.StringWidth(prefix) - ansi.StringWidth(label) - ansi.StringWidth(suffix)
	return style + prefix + label + strings.Repeat("─", max(fill, 0)) + suffix + terminalTone
}

func fitTerminalLabel(label string, width int) string {
	if width <= 0 {
		return ""
	}
	if ansi.StringWidth(label) <= width {
		return label
	}
	return ansi.Truncate(label, width, "…")
}

func terminalHeaderLabel(project string, shell string, width int) string {
	project = sanitizeTerminalLabel(project)
	shell = sanitizeTerminalLabel(shell)
	prefix := " phub / "
	separator := " · "
	shellLabel := "shell: " + shell + " "
	projectWidth := width - ansi.StringWidth(prefix) - ansi.StringWidth(separator) - ansi.StringWidth(shellLabel)
	if projectWidth < 1 {
		shellLabel = " " + shell + " "
		projectWidth = width - ansi.StringWidth(prefix) - ansi.StringWidth(separator) - ansi.StringWidth(shellLabel)
	}
	if projectWidth < 1 {
		return fitTerminalLabel(prefix+project+separator+shellLabel, width)
	}
	return prefix + fitTerminalLabel(project, projectWidth) + separator + shellLabel
}

func sanitizeTerminalLabel(value string) string {
	cleaned := ansi.Strip(value)
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, cleaned)
	return strings.Join(strings.Fields(cleaned), " ")
}

func terminalShellName(shell string) string {
	if shell == "" {
		return "shell"
	}
	return filepath.Base(shell)
}

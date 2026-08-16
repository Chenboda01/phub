package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const (
	ansiStyleReset      = "\x1b[0m"
	ansiStyleResetShort = "\x1b[m"
	ansiBackgroundReset = "\x1b[49m"
)

func (m model) paintBackground(content string) string {
	if m.width < 1 {
		return content
	}
	width := m.width

	fill := backgroundFill(m.theme.background)
	lines := strings.Split(content, "\n")
	var result strings.Builder
	for index, line := range lines {
		line = reapplyBackground(line, fill)
		result.WriteString(fill)
		result.WriteString(line)
		if padding := width - ansi.StringWidth(line); padding > 0 {
			result.WriteString(strings.Repeat(" ", padding))
		}
		if index < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func reapplyBackground(line string, fill string) string {
	line = strings.ReplaceAll(line, ansiStyleReset, ansiStyleReset+fill)
	line = strings.ReplaceAll(line, ansiStyleResetShort, ansiStyleResetShort+fill)
	line = strings.ReplaceAll(line, ansiBackgroundReset, ansiBackgroundReset+fill)
	return line
}

func backgroundFill(background color.Color) string {
	switch typed := background.(type) {
	case color.RGBA:
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", typed.R, typed.G, typed.B)
	default:
		red, green, blue, _ := background.RGBA()
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", red>>8, green>>8, blue>>8)
	}
}

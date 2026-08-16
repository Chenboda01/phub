package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestTerminalHeaderSanitizesUntrustedLabels(t *testing.T) {
	label := terminalHeaderLabel("Alpha\n\x1b[2J", "/bin/sh\r\x1b[31m", 37)

	if strings.ContainsAny(label, "\x1b\n\r") {
		t.Fatalf("header label contains control characters: %q", label)
	}
	if !strings.Contains(label, "Alpha") || !strings.Contains(label, "sh") {
		t.Fatalf("header label lost identity: %q", label)
	}
}

func TestTerminalHeaderFitsWideProjectNameAndKeepsShellContext(t *testing.T) {
	label := terminalHeaderLabel("東京プロジェクトの作業場", "fish", 37)
	line := terminalChromeLine(40, "╭─", label, "╮", terminalBold)
	visible := ansi.Strip(line)

	if width := ansi.StringWidth(visible); width != 40 {
		t.Fatalf("header width = %d, want 40: %q", width, visible)
	}
	if !strings.Contains(visible, "shell: fish") {
		t.Fatalf("header lost shell context: %q", visible)
	}
}

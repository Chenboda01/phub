package ui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestModelEntersThemeMenu_whenCtrlPIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))

	// Then
	if !next.themeMenu {
		t.Fatal("theme menu was not enabled")
	}
	if next.themeSelection != 0 {
		t.Fatalf("theme selection = %d, want 0", next.themeSelection)
	}
}

func TestModelKeepsThemePickerOpen_whenPrintableKeyIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	armed, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))
	next, _ := updateModel(t, armed, tea.KeyPressMsg(tea.Key{Text: "1", Code: '1'}))

	// Then
	if next.theme.name != defaultTheme.name {
		t.Fatalf("theme = %q, want %q", next.theme.name, defaultTheme.name)
	}
	if !next.themeMenu {
		t.Fatal("theme picker closed after printable key")
	}
}

func TestModelSelectsTheme_whenArrowNavigationThenEnterIsPressed(t *testing.T) {
	tests := []struct {
		name  string
		moves int
		want  string
		mode  themeMode
	}{
		{name: "red background", moves: 0, want: "Red Background", mode: backgroundMode},
		{name: "red theme", moves: 1, want: "Red Theme", mode: foregroundMode},
		{name: "red combo", moves: 2, want: "Red Combo", mode: comboMode},
		{name: "orange background", moves: 3, want: "Orange Background", mode: backgroundMode},
		{name: "orange theme", moves: 4, want: "Orange Theme", mode: foregroundMode},
		{name: "orange combo", moves: 5, want: "Orange Combo", mode: comboMode},
		{name: "yellow background", moves: 6, want: "Yellow Background", mode: backgroundMode},
		{name: "yellow theme", moves: 7, want: "Yellow Theme", mode: foregroundMode},
		{name: "yellow combo", moves: 8, want: "Yellow Combo", mode: comboMode},
		{name: "green background", moves: 9, want: "Green Background", mode: backgroundMode},
		{name: "green theme", moves: 10, want: "Green Theme", mode: foregroundMode},
		{name: "green combo", moves: 11, want: "Green Combo", mode: comboMode},
		{name: "blue background", moves: 12, want: "Blue Background", mode: backgroundMode},
		{name: "blue theme", moves: 13, want: "Blue Theme", mode: foregroundMode},
		{name: "blue combo", moves: 14, want: "Blue Combo", mode: comboMode},
		{name: "purple background", moves: 15, want: "Purple Background", mode: backgroundMode},
		{name: "purple theme", moves: 16, want: "Purple Theme", mode: foregroundMode},
		{name: "purple combo", moves: 17, want: "Purple Combo", mode: comboMode},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			current := newModel(testProjects(), "/bin/sh")

			// When
			armed, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))
			for move := 0; move < test.moves; move++ {
				armed, _ = updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
			}
			next, _ := updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

			// Then
			if next.theme.name != test.want {
				t.Fatalf("theme = %q, want %q", next.theme.name, test.want)
			}
			if next.theme.mode != test.mode {
				t.Fatalf("theme mode = %q, want %q", next.theme.mode, test.mode)
			}
			if next.themeMenu {
				t.Fatal("theme prefix remained armed after selection")
			}
		})
	}
}

func TestModelMovesThemeSelectionUp_whenUpArrowIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	armed, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))
	for move := 0; move < 5; move++ {
		armed, _ = updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	}

	// When
	for move := 0; move < 2; move++ {
		armed, _ = updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyUp}))
	}
	next, _ := updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	// Then
	if next.theme.name != "Orange Background" {
		t.Fatalf("theme = %q, want Orange Background", next.theme.name)
	}
}

func TestModelRendersThemeDropdown_whenCtrlPIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))
	content := next.render()

	// Then
	if !strings.Contains(content, "Choose a theme") {
		t.Fatalf("theme dropdown heading missing:\n%s", content)
	}
	if !strings.Contains(content, "> Red Background") {
		t.Fatalf("selected theme marker missing:\n%s", content)
	}
	if strings.Contains(content, "1 Red Background") {
		t.Fatalf("theme dropdown still requires numeric keys:\n%s", content)
	}
}

func TestModelOpensThemePickerOverTerminal_whenCtrlPIsPressed(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("still running")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: 'p', Mod: tea.ModCtrl}))

	// Then
	if !next.themeMenu {
		t.Fatal("Ctrl+P did not open the theme picker while the terminal is active")
	}
	if next.terminal == nil {
		t.Fatal("opening the theme picker closed the embedded terminal")
	}
	if len(session.keys) != 0 {
		t.Fatalf("Ctrl+P was forwarded to the shell: %v", session.keys)
	}
}

func TestModelRendersThemePickerOverTerminal_whenTerminalIsActive(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = &fakeTerminalSession{}
	current.themeMenu = true

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "Choose a theme") {
		t.Fatalf("theme dropdown heading missing while terminal is active:\n%s", content)
	}
}

func TestModelKeepsTerminalRunning_whenThemePickerIsCancelled(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("still running")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session
	current.themeMenu = true

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))

	// Then
	if next.themeMenu {
		t.Fatal("theme picker remained open after Esc")
	}
	if next.terminal == nil {
		t.Fatal("cancelling the theme picker closed the embedded terminal")
	}
}

func TestModelAppliesThemeAndKeepsTerminal_whenEnterIsPressed(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("still running")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session
	current.themeMenu = true

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	// Then
	if next.themeMenu {
		t.Fatal("theme picker remained open after Enter")
	}
	if next.terminal == nil {
		t.Fatal("applying a theme closed the embedded terminal")
	}
	if next.theme.name != "Red Background" {
		t.Fatalf("theme = %q, want Red Background", next.theme.name)
	}
}

func TestModelDoesNotForwardKeys_whenThemePickerIsOpenOverTerminal(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("still running")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session
	current.themeMenu = true

	// When
	_, _ = updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "e", Code: 'e'}))

	// Then
	if len(session.keys) != 0 {
		t.Fatalf("keys were forwarded to the shell while the theme picker was open: %v", session.keys)
	}
}

func TestModelPaintsFullWidthThemeBackground_whenViewIsRendered(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.width = 40
	current.height = 6

	// When
	content := current.paintBackground("phub\n\n> Alpha")

	// Then
	for _, line := range strings.Split(content, "\n") {
		if !strings.HasPrefix(line, "\x1b[48;2;") {
			t.Fatalf("line does not start with a background fill: %q", line)
		}
		if width := ansi.StringWidth(line); width != 40 {
			t.Fatalf("painted line width = %d, want 40: %q", width, line)
		}
	}
}

func TestModelReappliesBackgroundAfterStyleResets_whenLineContainsStyledRuns(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.width = 24

	// When
	content := current.paintBackground("\x1b[1mbold\x1b[0m then plain")

	// Then
	line := strings.Split(content, "\n")[0]
	if !strings.Contains(line, "\x1b[0m\x1b[48;2;") {
		t.Fatalf("background fill was not reapplied after style reset: %q", line)
	}
	if width := ansi.StringWidth(line); width != 24 {
		t.Fatalf("painted line width = %d, want 24: %q", width, line)
	}
}

func TestModelPadsWideCharactersToFullWidth_whenContentUsesCJK(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.width = 12

	// When
	content := current.paintBackground("東京")

	// Then
	line := strings.Split(content, "\n")[0]
	if width := ansi.StringWidth(line); width != 12 {
		t.Fatalf("painted CJK line width = %d, want 12: %q", width, line)
	}
}

func TestModelDoesNotPaintBackground_whenTerminalWidthIsUnknown(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	content := current.paintBackground("phub")

	// Then
	if content != "phub" {
		t.Fatalf("content = %q, want unchanged", content)
	}
}

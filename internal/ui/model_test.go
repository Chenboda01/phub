package ui

import (
	"context"
	"strings"
	"testing"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
)

func TestModelMovesSelection_whenNavigationKeyPressed(t *testing.T) {
	tests := []struct {
		name     string
		selected int
		key      tea.Key
		want     int
	}{
		{name: "j moves down", selected: 0, key: tea.Key{Text: "j", Code: 'j'}, want: 1},
		{name: "down arrow moves down", selected: 0, key: tea.Key{Code: tea.KeyDown}, want: 1},
		{name: "k moves up", selected: 1, key: tea.Key{Text: "k", Code: 'k'}, want: 0},
		{name: "up arrow moves up", selected: 1, key: tea.Key{Code: tea.KeyUp}, want: 0},
		{name: "down stops at final project", selected: 3, key: tea.Key{Text: "j", Code: 'j'}, want: 3},
		{name: "up stops at first project", selected: 0, key: tea.Key{Text: "k", Code: 'k'}, want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			current := newModel(testProjects(), "/bin/sh")
			current.selected = test.selected

			// When
			next, _ := updateModel(t, current, tea.KeyPressMsg(test.key))

			// Then
			if next.selected != test.want {
				t.Fatalf("selected = %d, want %d", next.selected, test.want)
			}
		})
	}
}

func TestModelEnablesSearchPlaceholder_whenSlashPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "/", Code: '/'}))

	// Then
	if !next.searching {
		t.Fatal("search placeholder was not enabled")
	}
}

func TestModelClearsSearchPlaceholder_whenEscapePressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.searching = true
	current.notice = "Search is coming soon."

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))

	// Then
	if next.searching {
		t.Fatal("search placeholder remained enabled")
	}
}

func TestModelStartsEmbeddedTerminal_whenEnterPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.startTerminal = func(context.Context, string, string, int, int) (terminalSession, error) {
		return nil, nil
	}

	// When
	next, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	// Then
	if command == nil {
		t.Fatal("Enter returned no terminal start command")
	}
	if !next.startingTerminal {
		t.Fatal("Enter did not enter terminal startup state")
	}
}

func TestModelStartsSelectedProjectTerminal_whenEnterPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.selected = 1
	var startedPath string
	current.startTerminal = func(_ context.Context, _ string, path string, _, _ int) (terminalSession, error) {
		startedPath = path
		return nil, nil
	}

	// When
	_, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	if command == nil {
		t.Fatal("Enter returned no terminal start command")
	}
	command()

	// Then
	if startedPath != testProjects()[1].Path {
		t.Fatalf("terminal project path = %q, want %q", startedPath, testProjects()[1].Path)
	}
}

func TestModelQuits_whenQPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	_, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}))

	// Then
	if command == nil {
		t.Fatal("quit command is nil")
	}
	result := command()
	if _, ok := result.(tea.QuitMsg); !ok {
		t.Fatalf("command result = %T, want tea.QuitMsg", result)
	}
}

func TestModelTracksTerminalSize_whenResized(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, tea.WindowSizeMsg{Width: 32, Height: 12})

	// Then
	if next.width != 32 || next.height != 12 {
		t.Fatalf("size = %dx%d, want 32x12", next.width, next.height)
	}
}

func TestModelRendersTerminalChrome_whenEmbeddedTerminalIsActive(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.width = 48
	current.height = 12
	current.terminal = &fakeTerminalSession{}
	current.terminalProject = "Alpha"

	// When
	content := current.render()

	// Then
	for _, expected := range []string{
		"phub / Alpha",
		"shell: sh",
		"Ctrl+D return",
		"Ctrl+C interrupt",
		"\x1b[1m",
		"\x1b[2m",
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("terminal chrome missing %q:\n%s", expected, content)
		}
	}
}

func TestModelOffsetsTerminalCursorBelowTerminalHeader(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = &fakeTerminalSession{cursorX: 7, cursorY: 3}

	// When
	view := current.View()

	// Then
	if view.Cursor == nil {
		t.Fatal("terminal cursor was not configured")
	}
	if view.Cursor.X != 7 || view.Cursor.Y != 4 {
		t.Fatalf("cursor = (%d, %d), want (7, 4)", view.Cursor.X, view.Cursor.Y)
	}
}

func TestModelKeepsSelectedProjectVisible_whenTerminalHeightShrinks(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.selected = 3
	current.height = 8

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "> Gamma") {
		t.Fatalf("rendered content does not show the selected project:\n%s", content)
	}
}

func TestModelAvoidsOverflow_whenTerminalIsNarrow(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.width = 40
	current.notice = "Opening Alpha is planned for a later milestone."

	// When
	content := current.render()

	// Then
	for _, line := range strings.Split(content, "\n") {
		if width := utf8.RuneCountInString(line); width > current.width {
			t.Fatalf("line width = %d, want at most %d: %q", width, current.width, line)
		}
	}
}

func TestModelShowsEmptyState_whenNoProjectsWereDiscovered(t *testing.T) {
	// Given
	current := newModel(nil, "/bin/sh")

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "No projects discovered.") {
		t.Fatalf("rendered content does not show the empty state:\n%s", content)
	}
}

func TestModelShowsNoOpenNotice_whenEnterPressedWithoutProjects(t *testing.T) {
	// Given
	current := newModel(nil, "/bin/sh")

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	// Then
	if next.notice != "No projects are available to open." {
		t.Fatalf("notice = %q", next.notice)
	}
}

func TestModelViewUsesOpaqueBackground_forEveryThemePreset(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When / Then
	for _, preset := range themes {
		current.theme = preset
		view := current.View()
		if view.BackgroundColor == nil {
			t.Fatalf("theme %q has no background color", preset.name)
		}
		foregroundSet := view.ForegroundColor != nil
		if wantForeground := preset.mode != backgroundMode; foregroundSet != wantForeground {
			t.Fatalf("theme %q foreground set = %t, want %t", preset.name, foregroundSet, wantForeground)
		}
	}
}

func TestModelShowsBuildVersion_whenProjectListIsRendered(t *testing.T) {
	// Given
	previous := Version
	Version = "abc1234"
	t.Cleanup(func() { Version = previous })
	current := newModel(testProjects(), "/bin/sh")

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "phub abc1234") {
		t.Fatalf("rendered content does not show the build version:\n%s", content)
	}
}

func TestModelShowsBuildVersion_whenCompactViewIsRendered(t *testing.T) {
	// Given
	previous := Version
	Version = "abc1234"
	t.Cleanup(func() { Version = previous })
	current := newModel(testProjects(), "/bin/sh")
	current.width = 60

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "phub abc1234") {
		t.Fatalf("compact view does not show the build version:\n%s", content)
	}
}

func TestModelOmitsBuildVersion_whenCompactFooterWouldOverflow(t *testing.T) {
	// Given
	previous := Version
	Version = "abc1234"
	t.Cleanup(func() { Version = previous })
	current := newModel(testProjects(), "/bin/sh")
	current.width = 40

	// When
	content := current.render()

	// Then
	if strings.Contains(content, "phub abc1234") {
		t.Fatalf("compact view shows the build version beyond the terminal width:\n%s", content)
	}
	for _, line := range strings.Split(content, "\n") {
		if width := utf8.RuneCountInString(line); width > current.width {
			t.Fatalf("line width = %d, want at most %d: %q", width, current.width, line)
		}
	}
}

func updateModel(t *testing.T, current model, message tea.Msg) (model, tea.Cmd) {
	t.Helper()

	next, command := current.Update(message)
	updated, ok := next.(model)
	if !ok {
		t.Fatalf("Update returned %T, want model", next)
	}

	return updated, command
}

func testProjects() []Project {
	return []Project{
		{Name: "Alpha", Path: "/tmp/Alpha"},
		{Name: "Delta", Path: "/tmp/Delta"},
		{Name: "Beta", Path: "/tmp/Beta"},
		{Name: "Gamma", Path: "/tmp/Gamma"},
	}
}

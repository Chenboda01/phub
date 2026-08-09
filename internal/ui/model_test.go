package ui

import (
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
			current := newModel()
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
	current := newModel()

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "/", Code: '/'}))

	// Then
	if !next.searching {
		t.Fatal("search placeholder was not enabled")
	}
}

func TestModelClearsSearchPlaceholder_whenEscapePressed(t *testing.T) {
	// Given
	current := newModel()
	current.searching = true
	current.notice = "Search is coming soon."

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))

	// Then
	if next.searching {
		t.Fatal("search placeholder remained enabled")
	}
}

func TestModelShowsOpenPlaceholder_whenEnterPressed(t *testing.T) {
	// Given
	current := newModel()

	// When
	next, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	// Then
	if next.notice != "Opening Forge is planned for a later milestone." {
		t.Fatalf("notice = %q", next.notice)
	}
}

func TestModelQuits_whenQPressed(t *testing.T) {
	// Given
	current := newModel()

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
	current := newModel()

	// When
	next, _ := updateModel(t, current, tea.WindowSizeMsg{Width: 32, Height: 12})

	// Then
	if next.width != 32 || next.height != 12 {
		t.Fatalf("size = %dx%d, want 32x12", next.width, next.height)
	}
}

func TestModelKeepsSelectedProjectVisible_whenTerminalHeightShrinks(t *testing.T) {
	// Given
	current := newModel()
	current.selected = 3
	current.height = 8

	// When
	content := current.render()

	// Then
	if !strings.Contains(content, "> Website") {
		t.Fatalf("rendered content does not show the selected project:\n%s", content)
	}
}

func TestModelAvoidsOverflow_whenTerminalIsNarrow(t *testing.T) {
	// Given
	current := newModel()
	current.width = 40
	current.notice = "Opening Forge is planned for a later milestone."

	// When
	content := current.render()

	// Then
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

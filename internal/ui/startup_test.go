package ui

import (
	"slices"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestStartupChooserLoadsGitHubProjectsByDefault_whenEnterIsPressed(t *testing.T) {
	// Given
	var requested ProjectScope
	current := newStartupModel("/bin/sh", func(scope ProjectScope) ([]Project, error) {
		requested = scope
		return []Project{{Name: "Alpha", Path: "/tmp/alpha"}}, nil
	})

	// When
	next, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	result := command()
	next, _ = updateModel(t, next, result)

	// Then
	if requested != ScopeGitHubOnly {
		t.Fatalf("requested scope = %d, want GitHub-only", requested)
	}
	if next.startup {
		t.Fatal("startup chooser remained open after loading")
	}
	if len(next.projects) != 1 || next.projects[0].Name != "Alpha" {
		t.Fatalf("projects = %q, want Alpha", next.projects)
	}
}

func TestStartupChooserLoadsAllLocalProjects_whenDownThenEnterIsPressed(t *testing.T) {
	// Given
	var requested ProjectScope
	current := newStartupModel("/bin/sh", func(scope ProjectScope) ([]Project, error) {
		requested = scope
		return []Project{{Name: "Local", Path: "/tmp/local"}}, nil
	})

	// When
	armed, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	next, command := updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	result := command()
	next, _ = updateModel(t, next, result)

	// Then
	if requested != ScopeAllLocal {
		t.Fatalf("requested scope = %d, want all-local", requested)
	}
	if next.startup {
		t.Fatal("startup chooser remained open after loading")
	}
}

func TestModelRefreshesUsingSelectedScope_whenRIsPressed(t *testing.T) {
	// Given
	requested := make([]ProjectScope, 0, 2)
	current := newStartupModel("/bin/sh", func(scope ProjectScope) ([]Project, error) {
		requested = append(requested, scope)
		return []Project{{Name: "Local", Path: "/tmp/local"}}, nil
	})
	armed, _ := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	loaded, command := updateModel(t, armed, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	loaded, _ = updateModel(t, loaded, command())

	// When
	refreshed, refreshCommand := updateModel(t, loaded, tea.KeyPressMsg(tea.Key{Text: "R", Code: 'r', Mod: tea.ModShift}))
	refreshed, _ = updateModel(t, refreshed, refreshCommand())

	// Then
	if !slices.Equal(requested, []ProjectScope{ScopeAllLocal, ScopeAllLocal}) {
		t.Fatalf("requested scopes = %v, want all-local for startup and refresh", requested)
	}
	if refreshed.startup {
		t.Fatal("startup chooser reopened during refresh")
	}
}

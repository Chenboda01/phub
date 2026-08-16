package ui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestModelRefreshesProjects_whenUppercaseRIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.selected = 3
	current.loadProjects = func(ProjectScope) ([]Project, error) {
		return []Project{{Name: "Alpha", Path: "/tmp/alpha"}}, nil
	}

	// When
	next, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "R", Code: 'r', Mod: tea.ModShift}))
	if command == nil {
		t.Fatal("R returned no refresh command")
	}
	if !next.refreshing {
		t.Fatal("refresh did not enter the refreshing state")
	}
	result := command()
	next, _ = updateModel(t, next, result)

	// Then
	if next.refreshing {
		t.Fatal("refresh remained in progress")
	}
	if len(next.projects) != 1 || next.projects[0].Path != "/tmp/alpha" {
		t.Fatalf("projects = %q, want refreshed Alpha project", next.projects)
	}
	if next.selected != 0 {
		t.Fatalf("selected = %d, want 0 after list shrinks", next.selected)
	}
	if next.notice != "Refreshed 1 projects." {
		t.Fatalf("notice = %q, want refresh summary", next.notice)
	}
}

func TestModelRefreshesProjects_whenLowercaseRIsPressed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.loadProjects = func(ProjectScope) ([]Project, error) {
		return []Project{{Name: "Alpha", Path: "/tmp/alpha"}}, nil
	}

	// When
	next, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "r", Code: 'r'}))

	// Then
	if command == nil {
		t.Fatal("r returned no refresh command")
	}
	if !next.refreshing {
		t.Fatal("refresh did not enter the refreshing state")
	}
}

func TestModelReportsRefreshError_whenUppercaseRScanFails(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.loadProjects = func(ProjectScope) ([]Project, error) {
		return nil, errors.New("scan failed")
	}

	// When
	_, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "R", Code: 'r', Mod: tea.ModShift}))
	result := command()
	next, _ := updateModel(t, current, result)

	// Then
	if next.refreshing {
		t.Fatal("refresh remained in progress after error")
	}
	if !strings.Contains(next.notice, "Could not refresh projects: scan failed") {
		t.Fatalf("notice = %q, want refresh error", next.notice)
	}
}

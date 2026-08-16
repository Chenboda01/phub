package ui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"phub/internal/git"
	"phub/internal/metadata"
)

func TestModelRendersMetadataColumns_whenMetadataLoaded(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.metadata = map[string]metadata.Info{
		"/tmp/Alpha": {Language: "Go", Git: git.Status{Branch: "main", Modified: 2}},
		"/tmp/Delta": {Language: "Python", Git: git.Status{Branch: "dev"}},
		"/tmp/Beta":  {Language: "Rust", Git: git.Status{Branch: "main", Untracked: 1}},
		"/tmp/Gamma": {Language: "JavaScript", Git: git.Status{Branch: "main"}},
	}

	// When
	content := current.render()

	// Then
	for _, expected := range []string{
		"Go", "main", "● 2 modified",
		"Python", "dev", "✓ clean",
		"Rust", "? 1 untracked",
		"JavaScript",
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("rendered content missing %q:\n%s", expected, content)
		}
	}
}

func TestModelLoadsMetadata_whenProjectsFinishLoading(t *testing.T) {
	// Given
	var loaded []string
	current := newModel(testProjects(), "/bin/sh")
	current.loadMetadata = func(_ context.Context, projects []Project) (map[string]metadata.Info, error) {
		for _, project := range projects {
			loaded = append(loaded, project.Path)
		}
		return map[string]metadata.Info{
			"/tmp/Alpha": {Language: "Go", Git: git.Status{Branch: "main"}},
		}, nil
	}

	// When
	next, command := updateModel(t, current, projectsRefreshedMsg{projects: testProjects()})
	if command == nil {
		t.Fatal("project load returned no metadata command")
	}
	result := command()
	next, _ = updateModel(t, next, result)

	// Then
	if len(loaded) != len(testProjects()) {
		t.Fatalf("metadata loader got %d projects, want %d", len(loaded), len(testProjects()))
	}
	if got := next.metadata["/tmp/Alpha"].Language; got != "Go" {
		t.Fatalf("language = %q, want Go", got)
	}
}

func TestModelShowsGitMissingNotice_whenMetadataLoaderReportsNotFound(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, projectsMetadataMsg{err: git.ErrNotFound})

	// Then
	if !strings.Contains(next.notice, "git executable not found") {
		t.Fatalf("notice = %q, want git missing warning", next.notice)
	}
}

func TestModelKeepsProjectsVisible_whenMetadataFails(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")

	// When
	next, _ := updateModel(t, current, projectsMetadataMsg{err: errors.New("metadata failed")})
	content := next.render()

	// Then
	if !strings.Contains(content, "Alpha") {
		t.Fatalf("projects disappeared after metadata failure:\n%s", content)
	}
	if !strings.Contains(next.notice, "Could not load project metadata") {
		t.Fatalf("notice = %q, want metadata error", next.notice)
	}
}

func TestModelRefreshesMetadata_whenEmbeddedTerminalExits(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("eof")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session
	current.terminalProject = "Alpha"
	loaded := false
	current.loadMetadata = func(_ context.Context, _ []Project) (map[string]metadata.Info, error) {
		loaded = true
		return nil, nil
	}

	// When
	next, command := updateModel(t, current, terminalOutputMsg{session: session, readErr: errors.New("eof")})
	if command == nil {
		t.Fatal("terminal exit returned no metadata refresh command")
	}
	next, _ = updateModel(t, next, command())

	// Then
	if !loaded {
		t.Fatal("metadata was not refreshed after terminal exit")
	}
	if next.terminal != nil {
		t.Fatal("embedded terminal remained active after exit")
	}
}

func TestModelRefreshesMetadata_whenProjectsAreRefreshed(t *testing.T) {
	// Given
	current := newModel(testProjects(), "/bin/sh")
	current.loadProjects = func(ProjectScope) ([]Project, error) {
		return []Project{{Name: "Alpha", Path: "/tmp/alpha"}}, nil
	}
	refreshed := false
	current.loadMetadata = func(_ context.Context, _ []Project) (map[string]metadata.Info, error) {
		refreshed = true
		return nil, nil
	}

	// When
	next, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "R", Code: 'r', Mod: tea.ModShift}))
	result := command()
	next, refreshCommand := updateModel(t, next, result)
	if refreshCommand == nil {
		t.Fatal("refresh returned no metadata command")
	}
	next, _ = updateModel(t, next, refreshCommand())

	// Then
	if !refreshed {
		t.Fatal("metadata was not refreshed after project refresh")
	}
}

package main

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"phub/internal/discovery"
)

func TestDiscoverProjects_returnsProjectsFromConfiguredRoot(t *testing.T) {
	// Given
	root := t.TempDir()
	projectDirectory := filepath.Join(root, "discovered-project")
	if err := os.Mkdir(projectDirectory, 0o700); err != nil {
		t.Fatalf("create project directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDirectory, "go.mod"), []byte{}, 0o600); err != nil {
		t.Fatalf("create project marker: %v", err)
	}
	t.Setenv(discovery.RootsEnvironment, root)
	t.Setenv(discovery.DepthEnvironment, "1")

	// When
	projects, err := discoverProjects(context.Background())

	// Then
	if err != nil {
		t.Fatalf("discoverProjects() error = %v", err)
	}
	if len(projects) != 1 || projects[0].Name != "discovered-project" {
		t.Fatalf("projects = %q, want discovered-project", projects)
	}
	if !slices.Equal([]string{projects[0].Path}, []string{projectDirectory}) {
		t.Fatalf("project path = %q, want %q", projects[0].Path, projectDirectory)
	}
}

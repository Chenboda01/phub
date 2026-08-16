package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"phub/internal/project"
)

func TestFilterProjects_keepsGitHubRemotes_whenGitHubOnlyIsSelected(t *testing.T) {
	// Given
	root := t.TempDir()
	githubPath := filepath.Join(root, "github-project")
	localPath := filepath.Join(root, "local-project")
	if err := os.Mkdir(githubPath, 0o700); err != nil {
		t.Fatalf("create GitHub project: %v", err)
	}
	if err := os.Mkdir(localPath, 0o700); err != nil {
		t.Fatalf("create local project: %v", err)
	}
	fakeGit := filepath.Join(t.TempDir(), "git")
	writeFakeGit(t, fakeGit, githubPath)
	t.Setenv("PATH", filepath.Dir(fakeGit))
	githubProject := newProject(t, githubPath)
	localProject := newProject(t, localPath)

	// When
	got, err := FilterProjects(context.Background(), []project.Project{githubProject, localProject}, GitHubOnly)

	// Then
	if err != nil {
		t.Fatalf("FilterProjects() error = %v", err)
	}
	if paths := projectPaths(got); !slices.Equal(paths, []string{githubPath}) {
		t.Fatalf("project paths = %q, want only GitHub project %q", paths, githubPath)
	}
}

func TestFilterProjects_keepsAllProjects_whenAllLocalIsSelected(t *testing.T) {
	// Given
	root := t.TempDir()
	firstPath := filepath.Join(root, "first")
	secondPath := filepath.Join(root, "second")
	if err := os.Mkdir(firstPath, 0o700); err != nil {
		t.Fatalf("create first project: %v", err)
	}
	if err := os.Mkdir(secondPath, 0o700); err != nil {
		t.Fatalf("create second project: %v", err)
	}
	projects := []project.Project{newProject(t, firstPath), newProject(t, secondPath)}

	// When
	got, err := FilterProjects(context.Background(), projects, AllLocal)

	// Then
	if err != nil {
		t.Fatalf("FilterProjects() error = %v", err)
	}
	if paths := projectPaths(got); !slices.Equal(paths, []string{firstPath, secondPath}) {
		t.Fatalf("project paths = %q, want %q", paths, []string{firstPath, secondPath})
	}
}

func TestIsGitHubRemoteURL_acceptsGitHubURLForms(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "https", url: "https://github.com/example/project.git", want: true},
		{name: "ssh scp", url: "git@github.com:example/project.git", want: true},
		{name: "ssh url", url: "ssh://git@github.com/example/project.git", want: true},
		{name: "other host", url: "https://gitlab.com/example/project.git", want: false},
		{name: "lookalike host", url: "https://github.com.example.com/project.git", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// When
			got := isGitHubRemoteURL(test.url)

			// Then
			if got != test.want {
				t.Fatalf("isGitHubRemoteURL(%q) = %t, want %t", test.url, got, test.want)
			}
		})
	}
}

func newProject(t *testing.T, path string) project.Project {
	t.Helper()
	result, err := project.NewFromDirectory(path)
	if err != nil {
		t.Fatalf("NewFromDirectory(%q): %v", path, err)
	}
	return result
}

func projectPaths(projects []project.Project) []string {
	paths := make([]string, len(projects))
	for index, current := range projects {
		paths[index] = current.Path()
	}
	return paths
}

func writeFakeGit(t *testing.T, path string, githubPath string) {
	t.Helper()
	script := fmt.Sprintf("#!/bin/sh\nif [ \"$2\" = %q ]; then\n  printf 'origin git@github.com:example/project.git (fetch)\\n'\nfi\n", githubPath)
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatalf("write fake git: %v", err)
	}
}

package project

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestHasMarker_reportsTrue_whenDirectoryContainsSupportedMarker(t *testing.T) {
	tests := []struct {
		name   string
		marker string
	}{
		{name: "git directory", marker: ".git"},
		{name: "Python project", marker: "pyproject.toml"},
		{name: "JavaScript project", marker: "package.json"},
		{name: "Rust project", marker: "Cargo.toml"},
		{name: "Go project", marker: "go.mod"},
		{name: "Java project", marker: "pom.xml"},
		{name: "Gradle project", marker: "build.gradle"},
		{name: "Kotlin Gradle project", marker: "build.gradle.kts"},
		{name: "CMake project", marker: "CMakeLists.txt"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			directory := t.TempDir()
			writeMarker(t, directory, test.marker)

			// When
			detected, err := HasMarker(directory)

			// Then
			if err != nil {
				t.Fatalf("HasMarker(%q) error = %v", directory, err)
			}
			if !detected {
				t.Fatalf("HasMarker(%q) = false, want true", directory)
			}
		})
	}
}

func TestHasMarker_reportsFalse_whenDirectoryHasNoSupportedMarker(t *testing.T) {
	// Given
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "notes.txt"), []byte("not a project"), 0o600); err != nil {
		t.Fatalf("write notes: %v", err)
	}

	// When
	detected, err := HasMarker(directory)

	// Then
	if err != nil {
		t.Fatalf("HasMarker(%q) error = %v", directory, err)
	}
	if detected {
		t.Fatalf("HasMarker(%q) = true, want false", directory)
	}
}

func TestNewFromDirectory_canonicalizesPath_whenDirectoryIsSymlink(t *testing.T) {
	// Given
	parent := t.TempDir()
	target := filepath.Join(parent, "actual-project")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatalf("create target: %v", err)
	}
	alias := filepath.Join(parent, "project-alias")
	if err := os.Symlink(target, alias); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	// When
	got, err := NewFromDirectory(alias)

	// Then
	if err != nil {
		t.Fatalf("NewFromDirectory(%q) error = %v", alias, err)
	}
	if got.Path() != target {
		t.Fatalf("Path() = %q, want %q", got.Path(), target)
	}
	if got.Name() != "actual-project" {
		t.Fatalf("Name() = %q, want %q", got.Name(), "actual-project")
	}
}

func TestNewFromDirectory_returnsNotDirectoryError_whenPathIsFile(t *testing.T) {
	// Given
	path := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// When
	_, err := NewFromDirectory(path)

	// Then
	if !errors.Is(err, ErrNotDirectory) {
		t.Fatalf("NewFromDirectory(%q) error = %v, want ErrNotDirectory", path, err)
	}
}

func writeMarker(t *testing.T, directory string, marker string) {
	t.Helper()

	path := filepath.Join(directory, marker)
	if marker == ".git" {
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatalf("create marker directory %q: %v", marker, err)
		}
		return
	}
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatalf("create marker file %q: %v", marker, err)
	}
}

package metadata

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"phub/internal/git"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name    string
		markers []string
		want    string
	}{
		{name: "go module", markers: []string{"go.mod"}, want: "Go"},
		{name: "rust crate", markers: []string{"Cargo.toml"}, want: "Rust"},
		{name: "javascript package", markers: []string{"package.json"}, want: "JavaScript"},
		{name: "typescript package", markers: []string{"package.json", "tsconfig.json"}, want: "TypeScript"},
		{name: "python project", markers: []string{"pyproject.toml"}, want: "Python"},
		{name: "python requirements", markers: []string{"requirements.txt"}, want: "Python"},
		{name: "java maven", markers: []string{"pom.xml"}, want: "Java"},
		{name: "java gradle", markers: []string{"build.gradle.kts"}, want: "Java"},
		{name: "c plus plus", markers: []string{"CMakeLists.txt"}, want: "C/C++"},
		{name: "unknown project", markers: nil, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			path := t.TempDir()
			for _, marker := range test.markers {
				if err := os.WriteFile(filepath.Join(path, marker), []byte(""), 0o600); err != nil {
					t.Fatalf("create marker %q: %v", marker, err)
				}
			}

			// When
			got := DetectLanguage(path)

			// Then
			if got != test.want {
				t.Fatalf("DetectLanguage() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestLoad_detectsLanguageWithoutGitRepository(t *testing.T) {
	// Given
	path := t.TempDir()
	if err := os.WriteFile(filepath.Join(path, "go.mod"), []byte(""), 0o600); err != nil {
		t.Fatalf("create go.mod: %v", err)
	}

	// When
	info := Load(context.Background(), path)

	// Then
	if info.Language != "Go" {
		t.Fatalf("language = %q, want Go", info.Language)
	}
	if info.GitErr != nil {
		t.Fatalf("GitErr = %v, want nil for non-repository", info.GitErr)
	}
}

func TestLoad_reportsMissingGitExecutable(t *testing.T) {
	// Given
	emptyPath := t.TempDir()
	t.Setenv("PATH", emptyPath)
	path := t.TempDir()
	if err := os.Mkdir(filepath.Join(path, ".git"), 0o700); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	// When
	info := Load(context.Background(), path)

	// Then
	if !errors.Is(info.GitErr, git.ErrNotFound) {
		t.Fatalf("GitErr = %v, want ErrNotFound", info.GitErr)
	}
}

func TestLoadAll_stopsOnCanceledContext(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// When
	_, err := LoadAll(ctx, []string{t.TempDir()})

	// Then
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestLoadAll_returnsInfoPerPath(t *testing.T) {
	// Given
	first := t.TempDir()
	second := t.TempDir()
	if err := os.WriteFile(filepath.Join(first, "pyproject.toml"), []byte(""), 0o600); err != nil {
		t.Fatalf("create pyproject.toml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(second, "package.json"), []byte(""), 0o600); err != nil {
		t.Fatalf("create package.json: %v", err)
	}

	// When
	infos, err := LoadAll(context.Background(), []string{first, second})

	// Then
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}
	if infos[first].Language != "Python" {
		t.Fatalf("first language = %q, want Python", infos[first].Language)
	}
	if infos[second].Language != "JavaScript" {
		t.Fatalf("second language = %q, want JavaScript", infos[second].Language)
	}
}

package discovery

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"phub/internal/project"
)

func TestScanner_Scan_discoversSortedProjects_whenRootsContainSupportedMarkers(t *testing.T) {
	// Given
	root := t.TempDir()
	createProject(t, root, "zeta", "go.mod")
	createProject(t, root, "alpha", "package.json")
	scanner := newScanner(t, 2)

	// When
	projects, err := scanner.Scan(context.Background(), []string{root})

	// Then
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if got := projectNames(projects); !slices.Equal(got, []string{"alpha", "zeta"}) {
		t.Fatalf("project names = %q, want %q", got, []string{"alpha", "zeta"})
	}
}

func TestScanner_Scan_respectsMaxDepth_whenProjectIsNestedBeyondLimit(t *testing.T) {
	// Given
	root := t.TempDir()
	shallow := createProject(t, root, "shallow", "go.mod")
	createProject(t, shallow, "too-deep", "Cargo.toml")
	scanner := newScanner(t, 1)

	// When
	projects, err := scanner.Scan(context.Background(), []string{root})

	// Then
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if got := projectNames(projects); !slices.Equal(got, []string{"shallow"}) {
		t.Fatalf("project names = %q, want %q", got, []string{"shallow"})
	}
}

func TestScanner_Scan_skipsIgnoredDirectories_whenNestedProjectHasMarker(t *testing.T) {
	// Given
	root := t.TempDir()
	ignored := filepath.Join(root, "node_modules")
	if err := os.Mkdir(ignored, 0o700); err != nil {
		t.Fatalf("create ignored directory: %v", err)
	}
	createProject(t, ignored, "hidden", "package.json")
	createProject(t, root, "visible", "go.mod")
	scanner := newScanner(t, 3)

	// When
	projects, err := scanner.Scan(context.Background(), []string{root})

	// Then
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if got := projectNames(projects); !slices.Equal(got, []string{"visible"}) {
		t.Fatalf("project names = %q, want %q", got, []string{"visible"})
	}
}

func TestScanner_Scan_doesNotModifyProjectFiles_whenProjectIsDiscovered(t *testing.T) {
	// Given
	root := t.TempDir()
	projectDirectory := createProject(t, root, "read-only", "go.mod")
	markerPath := filepath.Join(projectDirectory, "go.mod")
	beforeContents, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatalf("read marker before scan: %v", err)
	}
	beforeEntries, err := os.ReadDir(projectDirectory)
	if err != nil {
		t.Fatalf("read directory before scan: %v", err)
	}
	scanner := newScanner(t, 1)

	// When
	if _, err := scanner.Scan(context.Background(), []string{root}); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Then
	afterContents, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatalf("read marker after scan: %v", err)
	}
	afterEntries, err := os.ReadDir(projectDirectory)
	if err != nil {
		t.Fatalf("read directory after scan: %v", err)
	}
	if !bytes.Equal(afterContents, beforeContents) {
		t.Fatalf("marker contents changed from %q to %q", beforeContents, afterContents)
	}
	if len(afterEntries) != len(beforeEntries) || afterEntries[0].Name() != beforeEntries[0].Name() {
		t.Fatalf("directory entries changed from %v to %v", beforeEntries, afterEntries)
	}
}

func TestScanner_Scan_deduplicatesProjects_whenRootsOverlap(t *testing.T) {
	// Given
	root := t.TempDir()
	projectDirectory := createProject(t, root, "app", "go.mod")
	scanner := newScanner(t, 2)

	// When
	projects, err := scanner.Scan(context.Background(), []string{root, projectDirectory})

	// Then
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("project count = %d, want 1", len(projects))
	}
	if projects[0].Path() != projectDirectory {
		t.Fatalf("Path() = %q, want %q", projects[0].Path(), projectDirectory)
	}
}

func TestScanner_Scan_skipsMissingRoots_whenAnotherRootIsValid(t *testing.T) {
	// Given
	root := t.TempDir()
	createProject(t, root, "available", "go.mod")
	missing := filepath.Join(t.TempDir(), "missing")
	scanner := newScanner(t, 2)

	// When
	projects, err := scanner.Scan(context.Background(), []string{missing, root})

	// Then
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if got := projectNames(projects); !slices.Equal(got, []string{"available"}) {
		t.Fatalf("project names = %q, want %q", got, []string{"available"})
	}
}

func TestScanner_Scan_returnsCanceledError_whenContextIsCanceled(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	scanner := newScanner(t, 2)

	// When
	_, err := scanner.Scan(ctx, []string{t.TempDir()})

	// Then
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Scan() error = %v, want context.Canceled", err)
	}
}

func createProject(t *testing.T, parent string, name string, marker string) string {
	t.Helper()

	directory := filepath.Join(parent, name)
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatalf("create project directory %q: %v", name, err)
	}
	if err := os.WriteFile(filepath.Join(directory, marker), []byte{}, 0o600); err != nil {
		t.Fatalf("create marker %q: %v", marker, err)
	}

	return directory
}

func newScanner(t *testing.T, maxDepth int) Scanner {
	t.Helper()

	scanner, err := New(maxDepth)
	if err != nil {
		t.Fatalf("New(%d) error = %v", maxDepth, err)
	}

	return scanner
}

func projectNames(projects []project.Project) []string {
	names := make([]string, len(projects))
	for index, discovered := range projects {
		names[index] = discovered.Name()
	}

	return names
}

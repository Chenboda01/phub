package git

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReadStatus_reportsCleanRepository(t *testing.T) {
	// Given
	path := newRepository(t, "tracked.txt")
	runGit(t, path, "add", "tracked.txt")
	runGit(t, path, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "init")

	// When
	status, err := ReadStatus(context.Background(), path)

	// Then
	if err != nil {
		t.Fatalf("ReadStatus() error = %v", err)
	}
	if status.Dirty() {
		t.Fatalf("status = %+v, want clean", status)
	}
	if status.Branch != "master" && status.Branch != "main" {
		t.Fatalf("branch = %q, want master or main", status.Branch)
	}
}

func TestReadStatus_countsModifiedAndUntrackedFiles(t *testing.T) {
	// Given
	path := newRepository(t, "tracked.txt")
	runGit(t, path, "add", "tracked.txt")
	runGit(t, path, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "init")
	if err := os.WriteFile(filepath.Join(path, "tracked.txt"), []byte("changed"), 0o600); err != nil {
		t.Fatalf("modify tracked file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(path, "untracked.txt"), []byte("new"), 0o600); err != nil {
		t.Fatalf("create untracked file: %v", err)
	}

	// When
	status, err := ReadStatus(context.Background(), path)

	// Then
	if err != nil {
		t.Fatalf("ReadStatus() error = %v", err)
	}
	if status.Modified != 1 {
		t.Fatalf("modified = %d, want 1", status.Modified)
	}
	if status.Untracked != 1 {
		t.Fatalf("untracked = %d, want 1", status.Untracked)
	}
	if !status.Dirty() {
		t.Fatal("status = clean, want dirty")
	}
}

func TestReadStatus_reportsErrNotFound_whenGitIsMissing(t *testing.T) {
	// Given
	emptyPath := t.TempDir()
	t.Setenv("PATH", emptyPath)

	// When
	_, err := ReadStatus(context.Background(), t.TempDir())

	// Then
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}

func TestIsRepository_returnsFalse_whenDirectoryHasNoGit(t *testing.T) {
	// When
	got := IsRepository(t.TempDir())

	// Then
	if got {
		t.Fatal("IsRepository() = true, want false")
	}
}

func TestIsRepository_returnsTrue_whenDirectoryHasGit(t *testing.T) {
	// Given
	path := newRepository(t, "tracked.txt")

	// When
	got := IsRepository(path)

	// Then
	if !got {
		t.Fatal("IsRepository() = false, want true")
	}
}

func newRepository(t *testing.T, file string) string {
	t.Helper()
	path := t.TempDir()
	runGit(t, path, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(path, file), []byte("content"), 0o600); err != nil {
		t.Fatalf("create file: %v", err)
	}
	return path
}

func runGit(t *testing.T, path string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", path}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return string(output)
}

package launcher

import (
	"context"
	"errors"
	"testing"
)

func TestNewInteractiveShellCommandContext_startsShellInteractively(t *testing.T) {
	// Given
	projectPath := t.TempDir()

	// When
	command, err := NewInteractiveShellCommandContext(context.Background(), "/bin/sh", projectPath)

	// Then
	if err != nil {
		t.Fatalf("NewInteractiveShellCommandContext() error = %v", err)
	}
	if command.Dir != projectPath {
		t.Fatalf("Dir = %q, want %q", command.Dir, projectPath)
	}
	if len(command.Args) != 2 || command.Args[0] != "/bin/sh" || command.Args[1] != "-i" {
		t.Fatalf("Args = %q, want /bin/sh -i", command.Args)
	}
}

func TestNewShellCommand_setsExplicitWorkingDirectory_whenShellIsConfigured(t *testing.T) {
	// Given
	shell := "/bin/sh"
	projectPath := t.TempDir()

	// When
	command, err := NewShellCommand(shell, projectPath)

	// Then
	if err != nil {
		t.Fatalf("NewShellCommand() error = %v", err)
	}
	if command.Path != shell {
		t.Fatalf("Path = %q, want %q", command.Path, shell)
	}
	if command.Dir != projectPath {
		t.Fatalf("Dir = %q, want %q", command.Dir, projectPath)
	}
	if len(command.Args) != 1 {
		t.Fatalf("Args = %q, want only the executable", command.Args)
	}
}

func TestNewShellCommand_rejectsEmptyShell(t *testing.T) {
	// Given
	projectPath := t.TempDir()

	// When
	_, err := NewShellCommand("", projectPath)

	// Then
	if !errors.Is(err, ErrShellNotConfigured) {
		t.Fatalf("NewShellCommand() error = %v, want ErrShellNotConfigured", err)
	}
}

func TestNewShellCommand_rejectsEmptyProjectPath(t *testing.T) {
	// Given

	// When
	_, err := NewShellCommand("/bin/sh", "")

	// Then
	if !errors.Is(err, ErrProjectPathEmpty) {
		t.Fatalf("NewShellCommand() error = %v, want ErrProjectPathEmpty", err)
	}
}

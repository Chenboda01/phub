package launcher

import (
	"context"
	"errors"
	"os/exec"
)

var (
	ErrShellNotConfigured = errors.New("shell is not configured")
	ErrProjectPathEmpty   = errors.New("project path is empty")
)

func NewShellCommand(shell string, projectPath string) (*exec.Cmd, error) {
	if err := validateShellCommand(shell, projectPath); err != nil {
		return nil, err
	}

	command := exec.Command(shell)
	command.Dir = projectPath
	return command, nil
}

func NewInteractiveShellCommandContext(ctx context.Context, shell string, projectPath string) (*exec.Cmd, error) {
	if err := validateShellCommand(shell, projectPath); err != nil {
		return nil, err
	}

	command := exec.CommandContext(ctx, shell, "-i")
	command.Dir = projectPath
	return command, nil
}

func validateShellCommand(shell string, projectPath string) error {
	if shell == "" {
		return ErrShellNotConfigured
	}
	if projectPath == "" {
		return ErrProjectPathEmpty
	}
	return nil
}

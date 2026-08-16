package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrNotFound = errors.New("git executable not found")

type Status struct {
	Branch    string
	Modified  int
	Untracked int
}

func (s Status) Dirty() bool {
	return s.Modified > 0 || s.Untracked > 0
}

func IsRepository(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func ReadStatus(ctx context.Context, path string) (Status, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return Status{}, ErrNotFound
	}

	status := Status{Branch: readBranch(ctx, path)}
	output, err := exec.CommandContext(ctx, "git", "-C", path, "status", "--porcelain").Output()
	if err != nil {
		return Status{}, fmt.Errorf("read git status for %q: %w", path, err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "??"):
			status.Untracked++
		case len(line) >= 2 && (line[0] != ' ' || line[1] != ' '):
			status.Modified++
		}
	}

	return status, nil
}

func readBranch(ctx context.Context, path string) string {
	output, err := exec.CommandContext(ctx, "git", "-C", path, "symbolic-ref", "--quiet", "--short", "HEAD").Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	output, err = exec.CommandContext(ctx, "git", "-C", path, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

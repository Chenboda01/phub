package discovery

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"phub/internal/project"
)

type ProjectScope int

const (
	GitHubOnly ProjectScope = iota
	AllLocal
)

var ErrInvalidProjectScope = errors.New("invalid project scope")

func FilterProjects(ctx context.Context, projects []project.Project, scope ProjectScope) ([]project.Project, error) {
	switch scope {
	case AllLocal:
		return slices.Clone(projects), nil
	case GitHubOnly:
		filtered := make([]project.Project, 0, len(projects))
		for _, current := range projects {
			githubRemote, err := hasGitHubRemote(ctx, current.Path())
			if err != nil {
				return nil, err
			}
			if githubRemote {
				filtered = append(filtered, current)
			}
		}
		return filtered, nil
	default:
		return nil, ErrInvalidProjectScope
	}
}

func hasGitHubRemote(ctx context.Context, path string) (bool, error) {
	command := exec.CommandContext(ctx, "git", "-C", path, "remote", "-v")
	output, err := command.Output()
	if err != nil {
		if ctx.Err() != nil {
			return false, fmt.Errorf("check GitHub remote for %q: %w", path, ctx.Err())
		}
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return false, nil
		}
		return false, fmt.Errorf("read GitHub remote for %q: %w", path, err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && isGitHubRemoteURL(fields[1]) {
			return true, nil
		}
	}

	return false, nil
}

func isGitHubRemoteURL(remote string) bool {
	normalized := strings.ToLower(strings.TrimSpace(remote))
	for _, prefix := range []string{
		"https://github.com/",
		"https://www.github.com/",
		"http://github.com/",
		"http://www.github.com/",
		"git://github.com/",
		"ssh://git@github.com/",
		"git@github.com:",
	} {
		if strings.HasPrefix(normalized, prefix) {
			return true
		}
	}

	return false
}

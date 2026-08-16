package discovery

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"phub/internal/project"
)

var ErrInvalidMaxDepth = errors.New("maximum scan depth must be non-negative")

var ignoredDirectories = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	".venv":        {},
	"venv":         {},
	"target":       {},
	"dist":         {},
	"build":        {},
	"__pycache__":  {},
	".cache":       {},
	".cargo":       {},
	".config":      {},
	".local":       {},
	".npm":         {},
	".rustup":      {},
	"vendor":       {},
}

type Scanner struct {
	maxDepth int
}

func New(maxDepth int) (Scanner, error) {
	if maxDepth < 0 {
		return Scanner{}, ErrInvalidMaxDepth
	}

	return Scanner{maxDepth: maxDepth}, nil
}

func (s Scanner) Scan(ctx context.Context, roots []string) ([]project.Project, error) {
	found := make(map[string]project.Project)

	for _, root := range roots {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("scan projects: %w", err)
		}

		canonicalRoot, err := canonicalDirectory(root)
		if err != nil {
			if isSkippable(err) {
				continue
			}
			return nil, fmt.Errorf("resolve scan root %q: %w", root, err)
		}

		if err := s.scanRoot(ctx, canonicalRoot, found); err != nil {
			if isSkippable(err) {
				continue
			}
			return nil, fmt.Errorf("scan root %q: %w", canonicalRoot, err)
		}
	}

	projects := make([]project.Project, 0, len(found))
	for _, discovered := range found {
		projects = append(projects, discovered)
	}
	sort.Slice(projects, func(left int, right int) bool {
		if projects[left].Name() == projects[right].Name() {
			return projects[left].Path() < projects[right].Path()
		}

		return projects[left].Name() < projects[right].Name()
	})

	return projects, nil
}

func (s Scanner) scanRoot(ctx context.Context, root string, found map[string]project.Project) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if walkErr != nil {
			if isSkippable(walkErr) {
				if entry != nil && entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			return fmt.Errorf("inspect %q: %w", path, walkErr)
		}
		if !entry.IsDir() {
			return nil
		}
		if path != root && isIgnoredDirectory(entry.Name()) {
			return filepath.SkipDir
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("calculate scan depth for %q: %w", path, err)
		}
		depth := 0
		if relativePath != "." {
			depth = len(strings.Split(relativePath, string(filepath.Separator)))
		}
		if depth > s.maxDepth {
			return filepath.SkipDir
		}

		detected, err := project.HasMarker(path)
		if err != nil {
			if isSkippable(err) {
				return filepath.SkipDir
			}
			return fmt.Errorf("detect project markers in %q: %w", path, err)
		}
		if detected {
			discovered, err := project.NewFromDirectory(path)
			if err != nil {
				if isSkippable(err) {
					return filepath.SkipDir
				}
				return fmt.Errorf("create project from %q: %w", path, err)
			}
			found[discovered.Path()] = discovered
		}
		if depth == s.maxDepth {
			return filepath.SkipDir
		}

		return nil
	})
}

func canonicalDirectory(path string) (string, error) {
	canonicalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", path, err)
	}

	info, err := os.Stat(canonicalPath)
	if err != nil {
		return "", fmt.Errorf("inspect path %q: %w", canonicalPath, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("scan root %q: %w", canonicalPath, project.ErrNotDirectory)
	}

	return canonicalPath, nil
}

func isIgnoredDirectory(name string) bool {
	_, ok := ignoredDirectories[name]
	return ok
}

func isSkippable(err error) bool {
	return errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission)
}

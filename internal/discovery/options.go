package discovery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

const (
	RootsEnvironment = "PHUB_SCAN_ROOTS"
	DepthEnvironment = "PHUB_SCAN_DEPTH"
	defaultMaxDepth  = 4
)

var ErrInvalidDepth = errors.New("scan depth must be a non-negative integer")

type Options struct {
	roots    []string
	maxDepth int
}

func LoadOptions() (Options, error) {
	roots, err := loadRoots()
	if err != nil {
		return Options{}, err
	}

	maxDepth, err := loadMaxDepth()
	if err != nil {
		return Options{}, err
	}

	return Options{
		roots:    slices.Clone(roots),
		maxDepth: maxDepth,
	}, nil
}

func (o Options) Roots() []string {
	return slices.Clone(o.roots)
}

func (o Options) MaxDepth() int {
	return o.maxDepth
}

func loadRoots() ([]string, error) {
	rawRoots, configured := os.LookupEnv(RootsEnvironment)
	if configured {
		return nonEmptyPaths(filepath.SplitList(rawRoots)), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory for scan roots: %w", err)
	}

	return []string{
		home,
	}, nil
}

func loadMaxDepth() (int, error) {
	rawDepth, configured := os.LookupEnv(DepthEnvironment)
	if !configured || rawDepth == "" {
		return defaultMaxDepth, nil
	}

	maxDepth, err := strconv.Atoi(rawDepth)
	if err != nil || maxDepth < 0 {
		return 0, fmt.Errorf("%s %q: %w", DepthEnvironment, rawDepth, ErrInvalidDepth)
	}

	return maxDepth, nil
}

func nonEmptyPaths(paths []string) []string {
	filtered := make([]string, 0, len(paths))
	for _, path := range paths {
		if path != "" {
			filtered = append(filtered, path)
		}
	}

	return filtered
}

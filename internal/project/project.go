package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotDirectory = errors.New("project path is not a directory")

var markers = map[string]struct{}{
	".git":             {},
	"pyproject.toml":   {},
	"package.json":     {},
	"Cargo.toml":       {},
	"go.mod":           {},
	"pom.xml":          {},
	"build.gradle":     {},
	"build.gradle.kts": {},
	"CMakeLists.txt":   {},
}

type Project struct {
	path string
	name string
}

func NewFromDirectory(path string) (Project, error) {
	canonicalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Project{}, fmt.Errorf("resolve project path %q: %w", path, err)
	}

	info, err := os.Stat(canonicalPath)
	if err != nil {
		return Project{}, fmt.Errorf("inspect project path %q: %w", canonicalPath, err)
	}
	if !info.IsDir() {
		return Project{}, fmt.Errorf("project path %q: %w", canonicalPath, ErrNotDirectory)
	}

	return Project{
		path: canonicalPath,
		name: filepath.Base(canonicalPath),
	}, nil
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Path() string {
	return p.path
}

func HasMarker(directory string) (bool, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false, fmt.Errorf("read project directory %q: %w", directory, err)
	}

	for _, entry := range entries {
		if _, ok := markers[entry.Name()]; ok {
			return true, nil
		}
	}

	return false, nil
}

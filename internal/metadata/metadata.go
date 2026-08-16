package metadata

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"phub/internal/git"
)

type Info struct {
	Language string
	Git      git.Status
	GitErr   error
}

func Load(ctx context.Context, path string) Info {
	info := Info{Language: DetectLanguage(path)}
	if !git.IsRepository(path) {
		return info
	}

	status, err := git.ReadStatus(ctx, path)
	if err != nil {
		info.GitErr = err
		return info
	}
	info.Git = status

	return info
}

func LoadAll(ctx context.Context, paths []string) (map[string]Info, error) {
	result := make(map[string]Info, len(paths))
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return result, fmt.Errorf("load metadata: %w", err)
		}
		result[path] = Load(ctx, path)
	}
	return result, nil
}

func DetectLanguage(path string) string {
	switch {
	case fileExists(path, "go.mod"):
		return "Go"
	case fileExists(path, "Cargo.toml"):
		return "Rust"
	case fileExists(path, "package.json"):
		if fileExists(path, "tsconfig.json") {
			return "TypeScript"
		}
		return "JavaScript"
	case fileExists(path, "pyproject.toml"):
		return "Python"
	case fileExists(path, "requirements.txt"):
		return "Python"
	case fileExists(path, "setup.py"):
		return "Python"
	case fileExists(path, "Pipfile"):
		return "Python"
	case fileExists(path, "pom.xml"):
		return "Java"
	case fileExists(path, "build.gradle"):
		return "Java"
	case fileExists(path, "build.gradle.kts"):
		return "Java"
	case fileExists(path, "CMakeLists.txt"):
		return "C/C++"
	default:
		return ""
	}
}

func fileExists(directory string, name string) bool {
	_, err := os.Stat(filepath.Join(directory, name))
	return err == nil
}

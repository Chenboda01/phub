package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const sourceDirEnvironment = "PHUB_SOURCE_DIR"

func runAutoUpdate(stdout io.Writer, stderr io.Writer) int {
	updater, err := newAutoUpdater(stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "phub update failed: %v\n", err)
		return 1
	}
	code, err := updater.run()
	if err != nil {
		fmt.Fprintf(stderr, "phub update failed: %v\n", err)
		return 1
	}
	return code
}

type autoUpdater struct {
	stdout     io.Writer
	stderr     io.Writer
	sourceDir  string
	installDir string
	gitPath    string
	goPath     string
}

func newAutoUpdater(stdout io.Writer, stderr io.Writer) (autoUpdater, error) {
	sourceDir := os.Getenv(sourceDirEnvironment)
	if sourceDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return autoUpdater{}, fmt.Errorf("resolve home directory: %w", err)
		}
		sourceDir = filepath.Join(home, "phub")
	}

	executable, err := os.Executable()
	if err != nil {
		return autoUpdater{}, fmt.Errorf("resolve installed phub path: %w", err)
	}
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return autoUpdater{}, fmt.Errorf("git executable not found in $PATH: %w", err)
	}
	goPath, err := findGo()
	if err != nil {
		return autoUpdater{}, fmt.Errorf("go executable not found in $PATH; install Go or add its bin directory to your shell PATH (fish: fish_add_path <go-bin-dir>): %w", err)
	}

	return autoUpdater{
		stdout:     stdout,
		stderr:     stderr,
		sourceDir:  sourceDir,
		installDir: filepath.Dir(executable),
		gitPath:    gitPath,
		goPath:     goPath,
	}, nil
}

func findGo() (string, error) {
	if path, err := exec.LookPath("go"); err == nil {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err == nil {
		for _, directory := range []string{
			filepath.Join(home, ".local", "go", "bin"),
			filepath.Join(home, "go", "bin"),
		} {
			candidate := filepath.Join(directory, "go")
			if info, statErr := os.Stat(candidate); statErr == nil && !info.IsDir() {
				return candidate, nil
			}
		}
	}
	if info, err := os.Stat("/usr/local/go/bin/go"); err == nil && !info.IsDir() {
		return "/usr/local/go/bin/go", nil
	}

	return "", errors.New("go executable not found")
}

func (u autoUpdater) run() (int, error) {
	installPath := filepath.Join(u.installDir, "phub")

	fmt.Fprintf(u.stdout, "updating phub from %s\n", u.sourceDir)

	pull := exec.Command(u.gitPath, "pull")
	pull.Dir = u.sourceDir
	pull.Stdout = u.stdout
	pull.Stderr = u.stderr
	if err := pull.Run(); err != nil {
		return 0, fmt.Errorf("git pull: %w", err)
	}

	build := exec.Command(u.goPath, "build", "-o", installPath, "./cmd/phub")
	build.Dir = u.sourceDir
	build.Stdout = u.stdout
	build.Stderr = u.stderr
	if err := build.Run(); err != nil {
		return 0, fmt.Errorf("go build: %w", err)
	}
	fmt.Fprintf(u.stdout, "built %s\n", installPath)

	fmt.Fprintln(u.stdout, "starting phub")
	launch := exec.Command(installPath)
	launch.Stdin = os.Stdin
	launch.Stdout = u.stdout
	launch.Stderr = u.stderr
	if err := launch.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode(), nil
		}
		return 0, fmt.Errorf("start phub: %w", err)
	}

	return 0, nil
}

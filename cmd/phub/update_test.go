package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAutoUpdate_pullsBuildsAndLaunches(t *testing.T) {
	// Given
	binDir := fakeTools(t)
	sourceDir := t.TempDir()
	installDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "log")
	t.Setenv("QA_LOG", logPath)
	var stdout, stderr bytes.Buffer
	updater := autoUpdater{
		stdout:     &stdout,
		stderr:     &stderr,
		sourceDir:  sourceDir,
		installDir: installDir,
		gitPath:    filepath.Join(binDir, "git"),
		goPath:     filepath.Join(binDir, "go"),
	}

	// When
	code, err := updater.run()

	// Then
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("run() code = %d, want 0", code)
	}
	log, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("read tool log: %v", readErr)
	}
	content := string(log)
	if !strings.Contains(content, "git pull") {
		t.Fatalf("git was not asked to pull:\n%s", content)
	}
	if !strings.Contains(content, "go build -o "+filepath.Join(installDir, "phub")+" ./cmd/phub") {
		t.Fatalf("go was not asked to rebuild phub:\n%s", content)
	}
	if !strings.Contains(content, "dir "+sourceDir) {
		t.Fatalf("tools did not run inside the source directory:\n%s", content)
	}
	if !strings.Contains(content, "launched") {
		t.Fatalf("the rebuilt phub was not launched:\n%s", content)
	}
}

func TestAutoUpdate_stopsWithoutLaunching_whenPullFails(t *testing.T) {
	// Given
	binDir := fakeTools(t)
	logPath := filepath.Join(t.TempDir(), "log")
	t.Setenv("QA_LOG", logPath)
	t.Setenv("QA_GIT_FAIL", "1")
	var stdout, stderr bytes.Buffer
	updater := autoUpdater{
		stdout:     &stdout,
		stderr:     &stderr,
		sourceDir:  t.TempDir(),
		installDir: t.TempDir(),
		gitPath:    filepath.Join(binDir, "git"),
		goPath:     filepath.Join(binDir, "go"),
	}

	// When
	_, err := updater.run()

	// Then
	if err == nil || !strings.Contains(err.Error(), "git pull") {
		t.Fatalf("error = %v, want git pull failure", err)
	}
	log, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("read tool log: %v", readErr)
	}
	if strings.Contains(string(log), "launched") {
		t.Fatal("phub was launched after a failed pull")
	}
}

func TestAutoUpdate_stopsWithoutLaunching_whenBuildFails(t *testing.T) {
	// Given
	binDir := fakeTools(t)
	logPath := filepath.Join(t.TempDir(), "log")
	t.Setenv("QA_LOG", logPath)
	t.Setenv("QA_GO_FAIL", "1")
	var stdout, stderr bytes.Buffer
	updater := autoUpdater{
		stdout:     &stdout,
		stderr:     &stderr,
		sourceDir:  t.TempDir(),
		installDir: t.TempDir(),
		gitPath:    filepath.Join(binDir, "git"),
		goPath:     filepath.Join(binDir, "go"),
	}

	// When
	_, err := updater.run()

	// Then
	if err == nil || !strings.Contains(err.Error(), "go build") {
		t.Fatalf("error = %v, want go build failure", err)
	}
	log, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("read tool log: %v", readErr)
	}
	if strings.Contains(string(log), "launched") {
		t.Fatal("phub was launched after a failed build")
	}
}

func TestAutoUpdate_returnsLaunchedExitCode(t *testing.T) {
	// Given
	binDir := fakeTools(t)
	logPath := filepath.Join(t.TempDir(), "log")
	t.Setenv("QA_LOG", logPath)
	t.Setenv("QA_LAUNCH_EXIT", "1")
	var stdout, stderr bytes.Buffer
	updater := autoUpdater{
		stdout:     &stdout,
		stderr:     &stderr,
		sourceDir:  t.TempDir(),
		installDir: t.TempDir(),
		gitPath:    filepath.Join(binDir, "git"),
		goPath:     filepath.Join(binDir, "go"),
	}

	// When
	code, err := updater.run()

	// Then
	if err != nil {
		t.Fatalf("run() error = %v, want nil for launch exit", err)
	}
	if code != 3 {
		t.Fatalf("run() code = %d, want 3 from launched phub", code)
	}
}

func TestNewAutoUpdater_reportsMissingGit(t *testing.T) {
	// Given
	emptyPath := t.TempDir()
	t.Setenv("PATH", emptyPath)

	// When
	_, err := newAutoUpdater(&bytes.Buffer{}, &bytes.Buffer{})

	// Then
	if err == nil || !strings.Contains(err.Error(), "git executable not found") {
		t.Fatalf("error = %v, want missing git", err)
	}
}

func fakeTools(t *testing.T) string {
	t.Helper()
	binDir := t.TempDir()

	gitScript := `#!/bin/sh
echo "git $*" >> "$QA_LOG"
echo "dir $(pwd)" >> "$QA_LOG"
[ "$QA_GIT_FAIL" = "1" ] && exit 1
exit 0
`
	goScript := `#!/bin/sh
echo "go $*" >> "$QA_LOG"
echo "dir $(pwd)" >> "$QA_LOG"
prev=""
out=""
for arg in "$@"; do
  [ "$prev" = "-o" ] && out="$arg"
  prev="$arg"
done
[ "$QA_GO_FAIL" = "1" ] && exit 1
printf '#!/bin/sh\necho launched >> "$QA_LOG"\n[ "$QA_LAUNCH_EXIT" = "1" ] && exit 3\nexit 0\n' > "$out"
chmod +x "$out"
exit 0
`
	for name, script := range map[string]string{"git": gitScript, "go": goScript} {
		if err := os.WriteFile(filepath.Join(binDir, name), []byte(script), 0o700); err != nil {
			t.Fatalf("write fake %s: %v", name, err)
		}
	}

	return binDir
}

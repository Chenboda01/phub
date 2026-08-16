//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package terminal

import (
	"os/exec"
	"testing"
)

func TestConfigurePTYCommand_makesInteractiveShellsUseControllingTerminal(t *testing.T) {
	// Given
	command := exec.Command("/bin/sh", "-i")

	// When
	configurePTYCommand(command)

	// Then
	if command.SysProcAttr == nil {
		t.Fatal("PTY command has no process attributes")
	}
	if !command.SysProcAttr.Setsid {
		t.Fatal("PTY command does not create a session")
	}
	if !command.SysProcAttr.Setctty {
		t.Fatal("PTY command does not set a controlling terminal")
	}
}

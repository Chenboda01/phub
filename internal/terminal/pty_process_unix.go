//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package terminal

import (
	"os/exec"
	"syscall"
)

func configurePTYCommand(command *exec.Cmd) {
	if command.SysProcAttr == nil {
		command.SysProcAttr = &syscall.SysProcAttr{}
	}
	command.SysProcAttr.Setsid = true
	command.SysProcAttr.Setctty = true
}

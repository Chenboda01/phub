//go:build windows

package terminal

import "os/exec"

func configurePTYCommand(*exec.Cmd) {}

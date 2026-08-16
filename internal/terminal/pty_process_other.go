//go:build !windows && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !solaris

package terminal

import "os/exec"

func configurePTYCommand(*exec.Cmd) {}

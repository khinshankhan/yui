//go:build unix

package sysexec

import (
	"os/exec"
	"syscall"
)

func applyDetachedProcessAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

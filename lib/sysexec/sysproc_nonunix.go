//go:build !unix

package sysexec

import "os/exec"

func applyDetachedProcessAttrs(cmd *exec.Cmd) {
}

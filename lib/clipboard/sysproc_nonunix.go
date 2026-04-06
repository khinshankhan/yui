//go:build !unix

package clipboard

import "os/exec"

func applyDetachedProcessAttrs(cmd *exec.Cmd) {
}

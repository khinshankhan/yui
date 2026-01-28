package cobrawrapper

import (
	"strings"

	"github.com/spf13/cobra"
)

func localName(prefixCmds []string) string {
	if len(prefixCmds) == 0 {
		panic("must have at least 1 prefix cmd")
	}
	return prefixCmds[len(prefixCmds)-1]
}

func fullName(prefixCmds []string) string {
	if len(prefixCmds) == 0 {
		panic("must have at least 1 prefix cmd")
	}
	return strings.Join(prefixCmds, " ")
}

func CreateCmd(prefixCmds []string, cmd *cobra.Command) *cobra.Command {
	cmdName := localName(prefixCmds)
	fullCmdName := fullName(prefixCmds)

	actualCmd := *cmd
	actualCmd.Use = strings.ReplaceAll(cmd.Use, "%s", cmdName)
	actualCmd.Long = strings.ReplaceAll(cmd.Long, "%s", fullCmdName)

	return &actualCmd
}

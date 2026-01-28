package casecli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/caseconv"
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

func CreateCaseCmd(prefixCmds []string) *cobra.Command {
	cmdName := localName(prefixCmds)
	fullCmdName := fullName(prefixCmds)

	return &cobra.Command{
		Use: fmt.Sprintf(
			"%s [modes] [input]",
			cmdName,
		),
		Short: "Apply one or more casing transformations to input text",
		Long: fmt.Sprintf(
			`Transforms input using one or more case modes, like:
  %s lower kebab "Some Text"
  echo "Some Text" | %s kebab`,
			fullCmdName,
			fullCmdName,
		),
		Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			var input string
			var modes []string

			// Check if stdin is being piped
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				reader := bufio.NewReader(os.Stdin)
				in, _ := reader.ReadString('\n')
				input = strings.TrimSpace(in)
				modes = args
			} else {
				// Last arg is the input string, others are modes
				input = args[len(args)-1]
				modes = args[:len(args)-1]
			}

			for _, mode := range modes {
				input = caseconv.Convert(input, mode)
			}

			fmt.Println(input)
		},
	}
}

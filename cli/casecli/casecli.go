package casecli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/khinshankhan/yui/cli/cobrawrapper"
	"github.com/khinshankhan/yui/lib/caseconv"
	"github.com/spf13/cobra"
)

func CreateCaseCmd(prefixCmds []string) *cobra.Command {
	return cobrawrapper.CreateCmd(
		prefixCmds,
		&cobra.Command{
			Use:   "%s [modes] [input]",
			Short: "Apply one or more casing transformations to input text",
			Long: `Transforms input using one or more case modes, like:
  %s lower kebab "Some Text"
  echo "Some Text" | %s kebab`,
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
		},
	)
}

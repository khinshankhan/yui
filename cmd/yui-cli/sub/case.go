package sub

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var caseCmd = &cobra.Command{
	Use:   "case [modes] [input]",
	Short: "Apply one or more casing transformations to input text",
	Long: `Transforms input using one or more case modes, like:
  yui-cli case lower kebab "Some Text"\
  echo "Some Text" | yui-cli case kebab`,
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
			input = transform(input, mode)
		}

		fmt.Println(input)
	},
}

func splitWords(str string) []string {
	cleaned := strings.ToLower(strings.TrimSpace(str))
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	return strings.Fields(cleaned)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func transform(input, mode string) string {
	switch mode {
	case "upper":
		return strings.ToUpper(input)
	case "lower":
		return strings.ToLower(input)
	case "kebab":
		return strings.ToLower(strings.Join(splitWords(input), "-"))
	case "snake":
		return strings.ToLower(strings.Join(splitWords(input), "_"))
	case "pascal":
		parts := splitWords(input)
		for i, w := range parts {
			parts[i] = capitalize(w)
		}
		return strings.Join(parts, "")
	case "camel":
		parts := splitWords(input)
		if len(parts) == 0 {
			return ""
		}
		for i := range parts {
			if i == 0 {
				continue
			}
			parts[i] = capitalize(parts[i])
		}
		return parts[0] + strings.Join(parts[1:], "")
	default:
		return input
	}
}

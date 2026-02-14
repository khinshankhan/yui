package casecli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/caseconv"
)

func Usage(app string) string {
	return fmt.Sprintf(`Usage: %s [modes] [input]

Apply one or more casing transformations to input text.

Examples:
  %s lower kebab "Some Text"
  echo "Some Text" | %s kebab`, app, app, app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one argument is required")
	}

	var (
		input string
		modes []string
	)

	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		b, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			return fmt.Errorf("read stdin: %w", readErr)
		}
		input = strings.TrimSpace(string(b))
		modes = args
	} else {
		input = args[len(args)-1]
		modes = args[:len(args)-1]
	}

	for _, mode := range modes {
		input = caseconv.Convert(input, mode)
	}

	fmt.Println(input)
	return nil
}

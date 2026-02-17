package casecli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/caseconv"
)

var ErrHelpRequested = errors.New("help requested")

func isHelpArg(arg string) bool {
	switch arg {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func Help(app string) string {
	help := `%s - Text case conversion tools

USAGE:
    %s [flags] <text>
    echo "text" | %s [flags]

SUBCOMMANDS:
    help             Show this help message

CONVERSIONS:
    lower      Convert to lowercase
    upper      Convert to UPPERCASE
    kebab      Convert to kebab-case
    snake      Convert to snake_case
    camel      Convert to camelCase
    pascal     Convert to PascalCase

CHAINING:
    Chain multiple conversions by using multiple flags.
    Conversions are applied left-to-right.

EXAMPLES:
    %s lower "Hello World"             # hello world
    %s kebab "Hello World"             # hello-world
    %s snake upper "Hello World"       # HELLO_WORLD
    echo "Hello World" | %s kebab      # hello-world`

	return strings.ReplaceAll(help, "%s", app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one argument is required")
	}
	if len(args) == 1 && isHelpArg(args[0]) {
		return ErrHelpRequested
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

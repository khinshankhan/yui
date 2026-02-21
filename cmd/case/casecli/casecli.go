package casecli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/caseconv"
	"github.com/khinshankhan/yui/lib/cli"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Text case conversion tools").
		WithAliases(aliases...).
		WithArgs(cli.VariadicArg("conversion"), cli.RequiredArg("text")).
		WithSections(
			cli.Section{
				Title: "CONVERSIONS",
				Lines: []string{
					"lower      Convert to lowercase",
					"upper      Convert to UPPERCASE",
					"kebab      Convert to kebab-case",
					"snake      Convert to snake_case",
					"camel      Convert to camelCase",
					"pascal     Convert to PascalCase",
				},
			},
			cli.Section{
				Title: "TITLE CASE STYLES",
				Lines: []string{
					"apa        APA 7th Edition style",
					"chicago    Chicago Manual of Style 18th Edition",
					"mla        MLA Handbook 9th Edition",
					"ap         Associated Press 2020 Edition",
					"bluebook   Bluebook 21st Edition",
					"ama        AMA Manual of Style 11th Edition",
					"nytimes    NY Times style",
					"wikipedia  Wikipedia style",
				},
			},
			cli.Section{
				Title: "CHAINING",
				Lines: []string{
					"Chain multiple conversions by using multiple conversion tokens.",
					"Conversions are applied left-to-right.",
				},
			},
		).
		WithExamples(
			"%cmd% lower \"Hello World\"             # hello world",
			"%cmd% kebab \"Hello World\"             # hello-world",
			"%cmd% snake upper \"Hello World\"       # HELLO_WORLD",
			"echo \"Hello World\" | %cmd% kebab      # hello-world",
		).
		WithRun(run)
}

func run(ctx *cli.Context, args []string) error {
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

	fmt.Fprintln(ctx.Stdout, input)
	return nil
}

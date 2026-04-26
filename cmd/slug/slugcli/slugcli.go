package slugcli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/cli"
	"github.com/khinshankhan/yui/lib/slug"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Slug generation tools").
		WithAliases(aliases...).
		WithArgs(cli.OptionalArg("text")).
		RegisterFlags(
			cli.Flag{
				Name:        "reserved",
				Short:       "r",
				Value:       "slug",
				Description: "Reserve a slug value; repeat to allocate the next available match",
			},
		).
		WithExamples(
			"%cmd% \"Some Title\"                                # some-title",
			"%cmd% \"Some Title\" --reserved some-title          # some-title-1",
			"%cmd% \"Some Title\" -r some-title -r some-title-1  # some-title-2",
			"echo \"Some Title\" | %cmd%                         # some-title",
		).
		WithRun(run)
}

func run(ctx *cli.Context, args []string) error {
	text, reserved, err := parseArgs(args)
	if err != nil {
		return err
	}

	if text == "" {
		text, err = readStdin()
		if err != nil {
			return err
		}
	}

	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("text required via argument or stdin")
	}

	base := slug.Make(text)
	result := base
	if len(reserved) > 0 {
		result = slug.NextAvailable(base, reserved)
	}

	fmt.Fprintln(ctx.Stdout, result)
	return nil
}

func parseArgs(args []string) (string, []string, error) {
	var (
		textParts []string
		reserved  []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--reserved", "-r":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("%s requires a value", args[i])
			}
			reserved = append(reserved, args[i+1])
			i++
		default:
			if strings.HasPrefix(args[i], "--reserved=") {
				value := strings.TrimPrefix(args[i], "--reserved=")
				if value == "" {
					return "", nil, fmt.Errorf("--reserved requires a value")
				}
				reserved = append(reserved, value)
				continue
			}
			if strings.HasPrefix(args[i], "-") {
				return "", nil, fmt.Errorf("unknown flag: %s", args[i])
			}
			textParts = append(textParts, args[i])
		}
	}

	return strings.Join(textParts, " "), reserved, nil
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		b, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			return "", fmt.Errorf("read stdin: %w", readErr)
		}
		return strings.TrimSpace(string(b)), nil
	}

	return "", nil
}

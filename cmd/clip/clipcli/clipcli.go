package clipcli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/khinshankhan/yui/lib/cli"
	"github.com/khinshankhan/yui/lib/clipboard"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Clipboard tools").
		WithAliases(aliases...).
		WithSubcommandName("action").
		WithExample("Copy stdin into the clipboard", "copy").
		WithExample("Copy an argument into the clipboard", "copy", "\"hello world\"").
		WithExample("Paste clipboard contents", "paste").
		Register(
			NewCopyCommand("copy"),
			NewPasteCommand("paste"),
		)
}

func NewCopyCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Copy text into the clipboard").
		WithAliases(aliases...).
		WithArgs(cli.OptionalArg("text")).
		WithRun(runCopy)
}

func NewPasteCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Paste text from the clipboard").
		WithAliases(aliases...).
		WithRun(runPaste)
}

func runCopy(ctx *cli.Context, args []string) error {
	input, err := readClipboardInput(args)
	if err != nil {
		return err
	}

	if err := clipboard.Copy(input); err != nil {
		return err
	}

	return nil
}

func runPaste(ctx *cli.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("paste does not accept arguments")
	}

	text, err := clipboard.Paste()
	if err != nil {
		return err
	}

	_, err = io.WriteString(ctx.Stdout, text)
	return err
}

func readClipboardInput(args []string) (string, error) {
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		b, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			return "", fmt.Errorf("read stdin: %w", readErr)
		}
		return string(b), nil
	}

	if len(args) == 0 {
		return "", fmt.Errorf("text required via argument or stdin")
	}

	return strings.Join(args, " "), nil
}

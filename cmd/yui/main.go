package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/khinshankhan/yui/cmd/case/casecli"
	"github.com/khinshankhan/yui/cmd/color/colorcli"
	"github.com/khinshankhan/yui/cmd/net/netcli"
)

func Help(app string) string {
	help := `%s - A collection of micro tools

USAGE:
    %s [--format <format>] [--config <path>] [--db <path>] <command> [subcommand] [arguments]

COMMANDS:
    case, c      Text case conversion tools
    color, col   Color conversion tools
    net, n       Network/IP tools

GLOBAL COMMANDS:
    help         Show this help message

Use "%s <command> help" for more information about a command.`

	return strings.ReplaceAll(help, "%s", app)
}

func main() {
	help := Help("yui")

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	switch args[0] {
	case "case", "c":
		if err := casecli.Run(args[1:]); err != nil {
			if errors.Is(err, casecli.ErrHelpRequested) {
				fmt.Fprintln(os.Stdout, casecli.Help("yui case"))
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, casecli.Help("yui case"))
			os.Exit(1)
		}
	case "color", "col":
		if err := colorcli.Run(args[1:]); err != nil {
			if errors.Is(err, colorcli.ErrHelpRequested) {
				fmt.Fprintln(os.Stdout, colorcli.Help("yui color"))
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, colorcli.Help("yui color"))
			os.Exit(1)
		}
	case "net", "n":
		if err := netcli.Run(args[1:]); err != nil {
			if errors.Is(err, netcli.ErrHelpRequested) {
				fmt.Fprintln(os.Stdout, netcli.Help("yui net"))
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, netcli.Help("yui net"))
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Fprintln(os.Stdout, help)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	os.Exit(0)
}

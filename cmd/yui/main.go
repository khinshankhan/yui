package main

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/case/casecli"
	"github.com/khinshankhan/yui/cmd/color/colorcli"
	"github.com/khinshankhan/yui/cmd/net/netcli"
)

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: yui <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  case, c   Case conversion tools")
	fmt.Fprintln(os.Stderr, "  color     Color conversion tools")
	fmt.Fprintln(os.Stderr, "  net       Network tools")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "case", "c":
		if err := casecli.Run(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, casecli.Usage("yui case"))
			os.Exit(1)
		}
	case "color", "col":
		if err := colorcli.Run(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, colorcli.Usage("yui color"))
			os.Exit(1)
		}
	case "net", "n":
		if err := netcli.Run(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, netcli.Usage("yui net"))
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

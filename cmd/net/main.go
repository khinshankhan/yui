package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/net/netcli"
)

func main() {
	if err := netcli.Run(os.Args[1:]); err != nil {
		if errors.Is(err, netcli.ErrHelpRequested) {
			fmt.Fprintln(os.Stdout, netcli.Help("net"))
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, netcli.Help("net"))
		os.Exit(1)
	}

	os.Exit(0)
}

package main

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/net/netcli"
)

func main() {
	if err := netcli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, netcli.Usage("net"))
		os.Exit(1)
	}
}

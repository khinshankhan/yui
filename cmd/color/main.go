package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/color/colorcli"
)

func main() {
	if err := colorcli.Run(os.Args[1:]); err != nil {
		if errors.Is(err, colorcli.ErrHelpRequested) {
			fmt.Fprintln(os.Stdout, colorcli.Help("color"))
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, colorcli.Help("color"))
		os.Exit(1)
	}

	os.Exit(0)
}

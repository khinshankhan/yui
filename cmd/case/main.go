package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/case/casecli"
)

func main() {
	if err := casecli.Run(os.Args[1:]); err != nil {
		if errors.Is(err, casecli.ErrHelpRequested) {
			fmt.Fprintln(os.Stdout, casecli.Help("case"))
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, casecli.Help("case"))
		os.Exit(1)
	}

	os.Exit(0)
}

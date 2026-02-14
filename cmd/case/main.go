package main

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/case/casecli"
)

func main() {
	if err := casecli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, casecli.Usage("case"))
		os.Exit(1)
	}
}

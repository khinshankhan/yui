package main

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cmd/color/colorcli"
)

func main() {
	if err := colorcli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintln(os.Stderr, colorcli.Usage("color"))
		os.Exit(1)
	}
}

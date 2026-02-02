package main

import (
	"github.com/khinshankhan/yui/cli/colorcli"

	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(
		colorcli.
			CreateColorCmd([]string{"color"}).
			Execute(),
	)
}

func main() {
	Execute()
}

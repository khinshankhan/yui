package main

import (
	"github.com/khinshankhan/yui/cli/casecli"

	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(
		casecli.
			CreateCaseCmd([]string{"case"}).
			Execute(),
	)
}

func main() {
	Execute()
}

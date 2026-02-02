package main

import (
	"github.com/khinshankhan/yui/cli/casecli"
	"github.com/khinshankhan/yui/cli/colorcli"
	"github.com/khinshankhan/yui/cli/netcli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yui",
	Short: "A bespoke kit of tools — each built to solve one small problem very well.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(
		casecli.CreateCaseCmd([]string{"yui", "case"}),
		colorcli.CreateColorCmd([]string{"yui", "color"}),
		netcli.CreateNetCmd([]string{"yui", "net"}),
	)
}

func main() {
	Execute()
}

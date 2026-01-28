package main

import (
	"github.com/khinshankhan/yui/cli/netcli"

	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(
		netcli.
			CreateNetCmd([]string{"net"}).
			Execute(),
	)
}

func main() {
	Execute()
}

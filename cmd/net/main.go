package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/net/netcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := netcli.NewCommand("net")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

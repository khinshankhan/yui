package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/slug/slugcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := slugcli.NewCommand("slug")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

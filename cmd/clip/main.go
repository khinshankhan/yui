package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/clip/clipcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := clipcli.NewCommand("clip")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

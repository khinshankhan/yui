package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/sound/soundcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := soundcli.NewCommand("sound")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

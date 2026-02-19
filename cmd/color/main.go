package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/color/colorcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := colorcli.NewCommand("color")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

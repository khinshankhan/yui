package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/case/casecli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := casecli.NewCommand("case")
	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

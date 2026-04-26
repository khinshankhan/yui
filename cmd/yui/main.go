package main

import (
	"os"

	"github.com/khinshankhan/yui/cmd/case/casecli"
	"github.com/khinshankhan/yui/cmd/clip/clipcli"
	"github.com/khinshankhan/yui/cmd/color/colorcli"
	"github.com/khinshankhan/yui/cmd/net/netcli"
	"github.com/khinshankhan/yui/cmd/slug/slugcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func main() {
	root := cli.New("yui", "A collection of micro tools").
		WithSubcommandName("command").
		Register(
			casecli.NewCommand("case", "c"),
			slugcli.NewCommand("slug", "s"),
			clipcli.NewCommand("clip", "clipboard"),
			clipcli.NewCopyCommand("copy"),
			clipcli.NewPasteCommand("paste"),
			colorcli.NewCommand("color", "col"),
			netcli.NewCommand("net", "n"),
		)

	os.Exit(cli.Execute(root, os.Args[1:], os.Stdout, os.Stderr))
}

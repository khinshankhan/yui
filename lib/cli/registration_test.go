package cli_test

import (
	"testing"

	"github.com/khinshankhan/yui/cmd/case/casecli"
	"github.com/khinshankhan/yui/cmd/color/colorcli"
	"github.com/khinshankhan/yui/cmd/net/netcli"
	"github.com/khinshankhan/yui/lib/cli"
)

func TestCommandRegistrationsAreValid(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cli.Command
	}{
		{
			name: "case",
			cmd:  casecli.NewCommand("case", "c"),
		},
		{
			name: "color",
			cmd:  colorcli.NewCommand("color", "col"),
		},
		{
			name: "net",
			cmd:  netcli.NewCommand("net", "n"),
		},
		{
			name: "yui",
			cmd: cli.New("yui", "A collection of micro tools").
				WithSubcommandName("command").
				Register(
					casecli.NewCommand("case", "c"),
					colorcli.NewCommand("color", "col"),
					netcli.NewCommand("net", "n"),
				),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := cli.Validate(tc.cmd); err != nil {
				t.Fatalf("invalid command registration: %v", err)
			}
		})
	}
}

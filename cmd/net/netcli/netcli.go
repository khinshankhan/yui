package netcli

import (
	"fmt"

	"github.com/khinshankhan/yui/lib/cli"
	"github.com/khinshankhan/yui/lib/nettools"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	ip := cli.New("ip", "Get local IP addresses").
		WithAliases("i").
		WithSubcommandName("type").
		WithExample("Get primary local IP").
		WithExample("Get all interfaces and IPs", "all").
		WithDefaultSubcommand("primary").
		Register(
			cli.
				New("primary", "Get primary local IP").
				WithAliases("p").
				WithRun(runPrimary),
			cli.
				New("all", "Get all interfaces and IPs").
				WithAliases("a").
				WithRun(runAll),
		)

	return cli.New(name, "Network/IP tools").
		WithAliases(aliases...).
		WithExample("Show help", "help").
		WithExample("Get primary local IP", "ip").
		WithExample("Get all IPs", "ip", "all").
		Register(ip)
}

func runPrimary(ctx *cli.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown %s subcommand: %s", ctx.Command.Name, args[0])
	}

	ip, err := nettools.GetPrimaryLocalIP()
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Stdout, ip)
	return nil
}

func runAll(ctx *cli.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown %s subcommand: %s", ctx.Command.Name, args[0])
	}

	interfaces, err := nettools.GetLocalIPs()
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		fmt.Fprintf(ctx.Stdout, "%s:\n", iface.Name)
		if iface.MacAddress != "" {
			fmt.Fprintf(ctx.Stdout, "  MAC: %s\n", iface.MacAddress)
		}
		for _, ip := range iface.IPs {
			ipType := "IPv4"
			if ip.IsIPv6 {
				ipType = "IPv6"
			}
			fmt.Fprintf(ctx.Stdout, "  %s: %s\n", ipType, ip.Address)
		}
	}
	return nil
}

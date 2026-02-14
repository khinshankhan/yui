package netcli

import (
	"fmt"

	"github.com/khinshankhan/yui/lib/nettools"
)

func Usage(app string) string {
	return fmt.Sprintf(`Usage: %s <subcommand>

Subcommands:
  ip [primary|all]

Examples:
  %s ip
  %s ip primary
  %s ip all`, app, app, app, app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required")
	}

	switch args[0] {
	case "ip", "i":
		return runIP(args[1:])
	default:
		return fmt.Errorf("unknown net subcommand: %s", args[0])
	}
}

func runIP(args []string) error {
	subCmd := "primary"
	if len(args) > 0 {
		subCmd = args[0]
	}

	switch subCmd {
	case "primary", "p":
		ip, err := nettools.GetPrimaryLocalIP()
		if err != nil {
			return err
		}
		fmt.Println(ip)
		return nil
	case "all", "a":
		interfaces, err := nettools.GetLocalIPs()
		if err != nil {
			return err
		}

		for _, iface := range interfaces {
			fmt.Printf("%s:\n", iface.Name)
			if iface.MacAddress != "" {
				fmt.Printf("  MAC: %s\n", iface.MacAddress)
			}
			for _, ip := range iface.IPs {
				ipType := "IPv4"
				if ip.IsIPv6 {
					ipType = "IPv6"
				}
				fmt.Printf("  %s: %s\n", ipType, ip.Address)
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown ip subcommand: %s", subCmd)
	}
}

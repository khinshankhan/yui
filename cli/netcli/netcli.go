package netcli

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cli/cobrawrapper"
	"github.com/khinshankhan/yui/lib/nettools"
	"github.com/spf13/cobra"
)

func handleNetIP(args []string) {
	subCmd := "primary"
	if len(args) > 0 {
		subCmd = args[0]
	}

	switch subCmd {
	case "primary", "p":
		ip, err := nettools.GetPrimaryLocalIP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(ip)

	case "all", "a":
		interfaces, err := nettools.GetLocalIPs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
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

	default:
		fmt.Fprintf(os.Stderr, "Unknown ip subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

func CreateNetCmd(prefixCmds []string) *cobra.Command {
	return cobrawrapper.CreateCmd(
		prefixCmds,
		&cobra.Command{
			Use:   "%s <subcommand>",
			Short: "Suite of network tools.",
			Long: `Several network tools, like:
  %s ip         # Get primary local IP
  %s ip primary # Get primary local IP (default)
  %s ip all     # Get all IPs`,
			Args: cobra.MinimumNArgs(1),

			Run: func(cmd *cobra.Command, args []string) {
				subCmd := args[0]

				switch subCmd {
				case "ip", "i":
					handleNetIP(args[1:])
				default:
					fmt.Fprintf(os.Stderr, "Unknown net subcommand: %s\n", subCmd)
					os.Exit(1)
				}
			},
		},
	)
}

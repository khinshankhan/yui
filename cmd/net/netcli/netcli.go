package netcli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/khinshankhan/yui/lib/nettools"
)

var ErrHelpRequested = errors.New("help requested")

func isHelpArg(arg string) bool {
	switch arg {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func Help(app string) string {
	help := `%s - Network/IP tools

USAGE:
    %s <subcommand> [arguments]


SUBCOMMANDS:
    help               Show this help message

    ip [type]          Get local IP addresses
        primary, p     Get primary local IP (default)
        all, a         Get all interfaces and IPs

EXAMPLES:
    %s help                      # Show help
    %s ip                        # Get primary local IP
    %s ip all                    # Get all IPs

Use "%s <subcommand> help" for more information.`

	return strings.ReplaceAll(help, "%s", app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required")
	}
	if isHelpArg(args[0]) {
		return ErrHelpRequested
	}

	switch args[0] {
	case "ip", "i":
		return runIP(args[1:])
	default:
		return fmt.Errorf("unknown net subcommand: %s", args[0])
	}
}

func runIP(args []string) error {
	if len(args) > 0 && isHelpArg(args[0]) {
		return ErrHelpRequested
	}

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

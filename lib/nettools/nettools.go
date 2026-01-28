package nettools

import (
	"fmt"
	"net"
)

// NetworkInterface represents information about a network interface.
type NetworkInterface struct {
	Name       string   `json:"name"`
	MacAddress string   `json:"mac_address,omitempty"`
	IPs        []IPInfo `json:"ips,omitempty"`
}

// IPInfo represents IP address information.
type IPInfo struct {
	Address string `json:"address"`
	Network string `json:"network"`
	IsIPv6  bool   `json:"is_ipv6"`
}

// GetLocalIPs returns all local IP addresses from all network interfaces.
func GetLocalIPs() ([]NetworkInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	var result []NetworkInterface

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		netIface := NetworkInterface{
			Name:       iface.Name,
			MacAddress: iface.HardwareAddr.String(),
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var network string

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				network = v.Network()
			case *net.IPAddr:
				ip = v.IP
				network = v.Network()
			}

			if ip == nil {
				continue
			}

			// Skip loopback IPs
			if ip.IsLoopback() {
				continue
			}

			ipInfo := IPInfo{
				Address: ip.String(),
				Network: network,
				IsIPv6:  ip.To4() == nil,
			}

			netIface.IPs = append(netIface.IPs, ipInfo)
		}

		if len(netIface.IPs) > 0 {
			result = append(result, netIface)
		}
	}

	return result, nil
}

// getPrimaryFallback returns the first non-loopback IPv4 address found.
func getPrimaryFallback() (string, error) {
	interfaces, err := GetLocalIPs()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		for _, ip := range iface.IPs {
			if !ip.IsIPv6 {
				return ip.Address, nil
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}

// GetPrimaryLocalIP returns the primary local IP address (typically the one used for outbound connections).
func GetPrimaryLocalIP() (string, error) {
	// This trick finds the preferred outbound IP by connecting to an external address
	// (doesn't actually establish a connection, just determines the route)
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// Fallback: try to get any non-loopback IP
		return getPrimaryFallback()
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

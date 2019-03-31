package plugins

import (
	"net"
	"strings"
)

type IPs []string

func (ips IPs) String() string {
	out := ""
	if len(ips) > 0 {
		for i, ip := range ips {
			if i != 0 {
				out += ", "
			}
			out += ip
		}
	}
	return out
}

func GetIPs(interfaces ...string) (IPs, IPs, error) {
	var ipv4s, ipv6s []string
	for _, name := range interfaces {
		netIf, err := net.InterfaceByName(name)
		if err != nil {
			return nil, nil, err
		}
		addresses, err := netIf.Addrs()
		for _, addr := range addresses {
			ipString := addr.String()
			if strings.Contains(ipString, ":") {
				if !strings.HasPrefix(ipString, "fe80:") {
					ipv6s = append(ipv6s, ipString)
				}
			} else {
				ipv4s = append(ipv4s, ipString)
			}
		}
	}
	return IPs(ipv4s), IPs(ipv6s), nil
}

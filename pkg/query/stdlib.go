package query

import (
	"net"
)

type Stdlib struct {
	IPs map[string][]net.Addr
}

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

func (s *Stdlib) RefreshIPs(interfaces ...string) error {
	for _, name := range interfaces {
		netIf, err := net.InterfaceByName(name)
		if err != nil {
			return err
		}
		addresses, err := netIf.Addrs()
		if err != nil {
			return err
		}
		s.IPs[name] = addresses
	}
	return nil
}

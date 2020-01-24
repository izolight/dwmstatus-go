package dwmstatus

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type IPInfo struct {
	Interfaces []string
}

func (i *IPInfo) Refresh() string {
	ips := ""
	for _, name := range i.Interfaces {
		if name != "" {
			netIf, err := net.InterfaceByName(name)
			if err != nil {
				log.Println(err)
				return ""
			}
			addresses, err := netIf.Addrs()
			if err != nil {
				log.Println(err)
				return ""
			}
			for _, a := range addresses {
				if !strings.HasPrefix(a.String(), "fe80:") {
					if len(ips) > 0 {
						ips += ", "
					}
					ips = a.String()
				}
			}
		}
	}
	return fmt.Sprintf("IPs: %s", ips)
}

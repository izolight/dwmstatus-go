package dwmstatus

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type IPInfo struct {
	Interfaces []string
}

type ipUpdater struct {
	ipUpdate chan net.IPAddr
	ticker   *time.Ticker
}

func newIPUpdater(ipUpdate chan net.IPAddr, tickInterval int) ipUpdater {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println(conn.LocalAddr().(*net.UDPAddr).IP)

	return ipUpdater{
		ipUpdate: ipUpdate,
		ticker:   time.NewTicker(time.Duration(tickInterval) * time.Second),
	}
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

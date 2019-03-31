package main

import (
	"fmt"
	"log"
	"time"

	//"github.com/izolight/dwmstatus-go/pkg/dbus"
	//"github.com/izolight/dwmstatus-go/pkg/stdlib"

	//"github.com/izolight/dwmstatus-go/pkg/sysfs"
	"github.com/prometheus/procfs/sysfs"
)

type status string

func main() {
	interfaces := []string{"enp4s0", "enp5s0"}

	var status status
	for {
		status = ""
		/* ipv4s, ipv6s, err := stdlib.GetIPs(ifName, "enp0s31f6")
		if err != nil {
			log.Fatal(err)
		} else {
			if len(ipv4s) > 0 {
				status.addWithDelimiter("|", fmt.Sprintf("IP: %s", ipv4s))
			}
			if len(ipv6s) > 0 {
				status.addWithDelimiter("|", fmt.Sprintf("IPv6: %s", ipv6s))
			}
		}
		apPath, err := dbus.PathForAccessPoint(ifPath)
		if err != nil {
			log.Fatal(err)
		}
		ssid, err := dbus.SSID(apPath)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("SSID: %s", ssid))
		}
		bitrate, err := dbus.WifiLinkSpeed(ifPath)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("Speed: %s/s", humanize.Bytes(humanize.KByte*uint64(bitrate))))
		}
		rx, tx, err := sysfs.RxTxBytes(ifName)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("Down: %s/s Up: %s/s", humanize.Bytes(rx-prevRx), humanize.Bytes(tx-prevTx)))
		}
		prevRx = rx
		prevTx = tx */
		netClass, err := sysfs.NewNetClass()
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range interfaces {
			n := netClass[i]
			status.addWithDelimiter("|", fmt.Sprintf("%s: %s %d", i, n.OperState, *n.Speed))
		}

		fmt.Println(status)

		sleepUntil(1)
	}
}

func sleepUntil(seconds int) {
	sleepDuration := time.Duration(seconds) * time.Second
	now := time.Now()
	time.Sleep(now.Truncate(sleepDuration).Add(sleepDuration).Sub(now))
}

func (s *status) addWithDelimiter(delimiter string, value string) {
	if len(*s) != 0 {
		*s += status(" " + delimiter + " ")
	}
	*s += status(value)
}

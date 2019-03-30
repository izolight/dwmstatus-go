package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/godbus/dbus"
	"github.com/izolight/dwmstatus-go/pkg/sys"
	"github.com/izolight/dwmstatus-go/plugins"
)

type status string

func main() {
	ifName := "wlp4s0"
	prevRx, prevTx, err := sys.RxTxBytes(ifName)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	ifPath, err := plugins.GetDbusPathForInterface(ifName, conn)
	if err != nil {
		log.Fatal(err)
	}

	var status status
	for {
		status = ""
		ipv4s, ipv6s, err := plugins.GetIPs(ifName, "enp0s31f6")
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
		apPath, err := plugins.GetDbusPathForAP(ifPath, conn)
		if err != nil {
			log.Fatal(err)
		}
		ssid, err := plugins.GetSSIDFromDbus(apPath, conn)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("SSID: %s", ssid))
		}
		bitrate, err := plugins.GetBitrateFromDbus(ifPath, conn)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("Speed: %s/s", humanize.Bytes(humanize.KByte*uint64(bitrate))))
		}
		rx, tx, err := sys.RxTxBytes(ifName)
		if err == nil {
			status.addWithDelimiter("|", fmt.Sprintf("Down: %s/s Up: %s/s", humanize.Bytes(rx-prevRx), humanize.Bytes(tx-prevTx)))
		}
		prevRx = rx
		prevTx = tx
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
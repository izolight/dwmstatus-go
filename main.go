package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/godbus/dbus"
	"github.com/izolight/dwmstatus-go/plugins"
	"log"
	"time"
)

func printIPs(ipType string, ips[]string) string {
	if len(ips) > 0 {
		status := ipType + ": "
		for i, ip := range ips {
			if i != 0 {
				status += ", "
			}
			status += ip
			if i == len(ips)-1 {
				status += " | "
			}
		}
		return status
	}
	return ""
}

func main() {
	var rx, tx int
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	ifPath, err:= plugins.GetDbusPathForInterface("wlp4s0", conn)
	if err != nil {
		log.Fatal(err)
	}
	apPath, err := plugins.GetDbusPathForAP(ifPath, conn)
	if err != nil {
		log.Fatal(err)
	}

	for {
		status := ""
		ipv4s, ipv6s, err := plugins.GetIPs("wlp4s0", "enp0s31f6")
		if err != nil {
			status += fmt.Sprintf("Couldn't get ip addresses", err)
		} else {
			status += printIPs("IP", ipv4s)
			status += printIPs("IPv6", ipv6s)
		}
		wifiInfo, err := plugins.GetWifiInfo("wlp4s0")
		if err != nil {
			status += fmt.Sprintf("Couldn't get wifi info: %s", err)
		} else {
			status += wifiInfo.String()
			rxRate := humanize.Bytes(uint64(wifiInfo.RX - rx))
			txRate := humanize.Bytes(uint64(wifiInfo.TX - tx))
			status += fmt.Sprintf(" | Down: %s/s | Up: %s/s", rxRate, txRate)
			rx, tx = wifiInfo.RX, wifiInfo.TX
		}
		fmt.Println(status)
		ssid, err := plugins.GetSSIDFromDbus(apPath, conn)
		if err == nil {
			fmt.Println(ssid)
		}
		var now = time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	}
}

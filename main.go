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
		}
		return status
	}
	return ""
}

func main() {
	ifName := "wlp4s0"
	prevRx, prevTx, err := plugins.GetRxTxBytes(ifName)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	ifPath, err:= plugins.GetDbusPathForInterface(ifName, conn)
	if err != nil {
		log.Fatal(err)
	}
	apPath, err := plugins.GetDbusPathForAP(ifPath, conn)
	if err != nil {
		log.Fatal(err)
	}

	var status string
	for {
		status = ""
		ipv4s, ipv6s, err := plugins.GetIPs(ifName, "enp0s31f6")
		if err != nil {
			status += fmt.Sprintf("Couldn't get ip addresses", err)
		} else {
			status += printIPs("IP", ipv4s)
			status += printIPs(" | IPv6", ipv6s)
		}
		ssid, err := plugins.GetSSIDFromDbus(apPath, conn)
		if err == nil {
			status += fmt.Sprintf(" | SSID: %s", ssid)
		}
		bitrate, err := plugins.GetBitrateFromDbus(ifPath, conn)
		if err == nil {
			status += fmt.Sprintf(" | Speed: %s/s", humanize.Bytes(humanize.KByte *uint64(bitrate)))
		}
		rx, tx, err := plugins.GetRxTxBytes(ifName)
		if err == nil {
			status += fmt.Sprintf(" | Down: %s/s | Up: %s/s", humanize.Bytes(rx-prevRx), humanize.Bytes(tx-prevTx))
		}
		prevRx = rx
		prevTx = tx
/*		wifiInfo, err := plugins.GetWifiInfo("wlp4s0")
		if err != nil {
			status += fmt.Sprintf("Couldn't get wifi info: %s", err)
		} else {
			status += wifiInfo.String()
			rxRate := humanize.Bytes(uint64(wifiInfo.RX - rx))
			txRate := humanize.Bytes(uint64(wifiInfo.TX - tx))
			status += fmt.Sprintf(" | Down: %s/s | Up: %s/s", rxRate, txRate)
			rx, tx = wifiInfo.RX, wifiInfo.TX
		}*/
		fmt.Println(status)

		sleepUntil(5)
	}
}

func sleepUntil(seconds int) {
	sleepDuration := time.Duration(seconds) * time.Second
	now := time.Now()
	time.Sleep(now.Truncate(sleepDuration).Add(sleepDuration).Sub(now))
}
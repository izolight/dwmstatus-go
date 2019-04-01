package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/mdlayher/wifi"
	"github.com/prometheus/procfs"
)

var config struct {
	wifiInterface string
	lanInterfaces []string
	vpnInterfaces []string
	allInterfaces []string
	batteries     []string
}

type status string

func main() {
	config.wifiInterface = "wlp4s0"
	config.lanInterfaces = []string{"enp0s31f6"}
	config.vpnInterfaces = []string{"home", "remote"}

	config.allInterfaces = append(config.allInterfaces, config.wifiInterface)
	config.allInterfaces = append(config.allInterfaces, config.lanInterfaces...)
	config.allInterfaces = append(config.allInterfaces, config.vpnInterfaces...)

	wifiClient, err := wifi.New()
	if err != nil {
		log.Fatal(err)
	}
	interfaces, err := wifiClient.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	wifiInterface := &wifi.Interface{}
	for _, ifi := range interfaces {
		if ifi.Name == config.wifiInterface {
			wifiInterface = ifi
		}
	}

	var status status
	var prevRx, prevTx uint64
	for {
		status = ""
		bss, err := wifiClient.BSS(wifiInterface)
		if err != nil {
			log.Println(err)
		}
		stationInfo, err := wifiClient.StationInfo(wifiInterface)
		if err != nil {
			log.Println(err)
		}
		status.addWithDelimiter("|", fmt.Sprintf("SSID: %s %d%%(%ddBm)", bss.SSID, wifiPercentage(stationInfo[0].Signal), stationInfo[0].Signal))

		nd, err := procfs.NewNetDev()
		if err != nil {
			log.Panicln(err)
		}

		total := nd.Total()
		rx, tx := total.RxBytes, total.TxBytes
		status.addWithDelimiter("|", fmt.Sprintf("U:%s/s", humanize.Bytes(tx-prevTx)))
		status.addWithDelimiter("|", fmt.Sprintf("D:%s/s", humanize.Bytes(rx-prevRx)))
		prevRx, prevTx = rx, tx

		fmt.Println(status)
		sleepUntil(1)
	}
	/*

		var status status
		for {
			status = ""
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
		} */
}

func wifiPercentage(signal int) int {
	if signal <= -100 {
		return 0
	} else if signal >= -50 {
		return 100
	}
	return 2 * (signal + 100)
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

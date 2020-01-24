package main

import (
	"flag"
	"fmt"
	"github.com/izolight/dwmstatus-go/pkg"
	"log"
	"strings"
	"time"
)

type status string

func main() {
	wifiInterface := flag.String("wifi", "", "")
	lanInterfaces := flag.String("lan", "", "")
	batteries := flag.String("bat", "", "")
	flag.Parse()

	// Add Plugins
	transferInfo := &dwmstatusgo.TransferInfo{}
	ipInfo := &dwmstatusgo.IPInfo{
		Interfaces: strings.Split(*lanInterfaces, ","),
	}
	plugins := []dwmstatusgo.Refresher{
		transferInfo,
		ipInfo,
	}
	if len(*wifiInterface) != 0 {
		wifiStats := &dwmstatusgo.WifiStats{
			IfiName: *wifiInterface,
		}
		err := wifiStats.Initialise()
		if err != nil {
			log.Println(err)
		}
		ipInfo.Interfaces = append(ipInfo.Interfaces, *wifiInterface)
		plugins = append(plugins, wifiStats)
	}
	if len(*batteries) != 0 {
		batteryStats := &dwmstatusgo.BatteryStats{
			Batteries:strings.Split(*batteries, ","),
		}
		plugins = append(plugins, batteryStats)
	}
	cpuInfo := &dwmstatusgo.CPUInfo{}
	memInfo := dwmstatusgo.NewMemoryInfo()
	plugins = append(plugins, cpuInfo, memInfo)

	var status status
	for {
		status = ""
		for _, r := range plugins {
			status.addWithDelimiter("|", r.Refresh())
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
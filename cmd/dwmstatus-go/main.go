package main

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/procfs/sysfs"

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
	config.batteries = []string{"BAT0", "BAT1"}

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
	var prevEnergy int64
	var remaining int64
	nextMeasurement := time.Now().Truncate(time.Second).Add(5 * time.Second)
	for {
		status = ""

		// wifi stuff
		bss, err := wifiClient.BSS(wifiInterface)
		if err != nil {
			log.Println(err)
		}
		stationInfo, err := wifiClient.StationInfo(wifiInterface)
		if err != nil {
			log.Println(err)
		}
		status.addWithDelimiter("|", fmt.Sprintf("SSID: %s %d%%(%ddBm)", bss.SSID, wifiPercentage(stationInfo[0].Signal), stationInfo[0].Signal))

		// transfer stats
		nd, err := procfs.NewNetDev()
		if err != nil {
			log.Println(err)
		}
		total := nd.Total()
		rx, tx := total.RxBytes, total.TxBytes
		status.addWithDelimiter("|", fmt.Sprintf("U:%s/s", humanize.Bytes(tx-prevTx)))
		status.addWithDelimiter("|", fmt.Sprintf("D:%s/s", humanize.Bytes(rx-prevRx)))
		prevRx, prevTx = rx, tx

		// battery stats
		psc, err := sysfs.NewPowerSupplyClass()
		if err != nil {
			log.Println(err)
		}
		var energyNow, energyFull, EnergyFullDesign int64
		for _, b := range config.batteries {
			bat, ok := psc[b]
			if !ok {
				continue
			}
			energyNow += *bat.EnergyNow
			energyFull += *bat.EnergyFull
			EnergyFullDesign += *bat.EnergyFullDesign
		}
		measure := time.Now().Truncate(time.Second)
		if measure.UnixNano() == nextMeasurement.UnixNano() {
			nextMeasurement = nextMeasurement.Add(5 * time.Second)
			remaining = calculateRemainingTime(energyNow, prevEnergy, energyFull)
			prevEnergy = energyNow
		}
		status.addWithDelimiter("|", fmt.Sprintf("BAT: %.1f%% %dMin", (float64(energyNow)/float64(energyFull))*100, remaining))

		fmt.Println(status)
		sleepUntil(1)
	}
}

func wifiPercentage(signal int) int {
	if signal <= -100 {
		return 0
	} else if signal >= -50 {
		return 100
	}
	return 2 * (signal + 100)
}

func calculateRemainingTime(energyNow int64, energyPrev int64, EnergyFull int64) int64 {
	remainingEnergy := EnergyFull - energyNow
	charged := energyNow - energyPrev
	if charged == 0 {
		return 0
	} else if charged > 0 {
		minutes := float64(remainingEnergy) / float64(charged) / 60 * 5
		return int64(minutes)
	} else {
		minutes := float64(energyNow) / float64(charged) / 60 * 5
		return int64(minutes)
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

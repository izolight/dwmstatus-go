package main

import (
	"fmt"
	"log"
	"net"
	"strings"
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

var info struct {
	status          status
	prevRx, prevTx  uint64
	prevEnergy      int64
	remaining       int64
	nextMeasurement time.Time
	wifiInterface   *wifi.Interface
	wifiClient      *wifi.Client
	bss             *wifi.BSS
	stationInfo     *wifi.StationInfo
}

type status string

func main() {
	config.wifiInterface = "wlp4s0"
	config.lanInterfaces = []string{"enp0s31f6"}
	config.vpnInterfaces = []string{"home", "remote"}
	config.batteries = []string{"BAT0", "BAT1"}

	config.allInterfaces = append(config.allInterfaces, config.wifiInterface)
	config.allInterfaces = append(config.allInterfaces, config.lanInterfaces...)

	err := initialise()
	if err != nil {
		log.Fatal(err)
	}

	for {
		info.status = ""

		ips, err := refreshIPs(config.allInterfaces)
		if err != nil {
			log.Println(err)
		}
		info.status.addWithDelimiter("|", fmt.Sprintf("IPs: %s", ips))

		// wifi stuff
		bss, err := info.wifiClient.BSS(info.wifiInterface)
		if err != nil {
			log.Println(err)
		}
		stationInfo, err := info.wifiClient.StationInfo(info.wifiInterface)
		if err != nil {
			log.Println(err)
		}
		info.status.addWithDelimiter("|", fmt.Sprintf("SSID: %s %d%%(%ddBm)", bss.SSID, wifiPercentage(stationInfo[0].Signal), stationInfo[0].Signal))

		// transfer stats
		nd, err := procfs.NewNetDev()
		if err != nil {
			log.Println(err)
		}
		total := nd.Total()
		rx, tx := total.RxBytes, total.TxBytes
		info.status.addWithDelimiter("|", fmt.Sprintf("U:%s/s", humanize.Bytes(tx-info.prevTx)))
		info.status.addWithDelimiter("|", fmt.Sprintf("D:%s/s", humanize.Bytes(rx-info.prevRx)))
		info.prevRx, info.prevTx = rx, tx

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
		if measure.UnixNano() == info.nextMeasurement.UnixNano() {
			info.nextMeasurement = info.nextMeasurement.Add(5 * time.Second)
			info.remaining = calculateRemainingTime(energyNow, info.prevEnergy, energyFull)
			info.prevEnergy = energyNow
		}
		info.status.addWithDelimiter("|", fmt.Sprintf("BAT: %.1f%% %dMin", (float64(energyNow)/float64(energyFull))*100, info.remaining))

		fmt.Println(info.status)
		sleepUntil(1)
	}
}

func initialise() error {
	var err error
	info.wifiClient, err = wifi.New()
	if err != nil {
		return err
	}
	interfaces, err := info.wifiClient.Interfaces()
	if err != nil {
		return err
	}
	for _, ifi := range interfaces {
		if ifi.Name == config.wifiInterface {
			info.wifiInterface = ifi
		}
	}
	info.nextMeasurement = time.Now().Truncate(time.Second).Add(5 * time.Second)
	return nil
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
		return int64(minutes) * (-1)
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

func refreshIPs(interfaces []string) (string, error) {
	ips := ""
	for _, name := range interfaces {
		netIf, err := net.InterfaceByName(name)
		if err != nil {
			return "", err
		}
		addresses, err := netIf.Addrs()
		if err != nil {
			return "", err
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
	return ips, nil
}

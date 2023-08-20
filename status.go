package dwmstatus

import "net"

import "fmt"

import "time"

const statusBar string = "%s | IP: %s | SSID: %s %d%% | BAT: %d%% | CPU: %d%% | RAM: %d | %s"

type status struct {
	lastBatteryPercentage   uint
	lastRAMUsage            uint
	lastCPUUsage            uint
	lastWifi                wifiStatus
	lastIP                  net.IPAddr
	lastTransferInfo        transferInfo
	batteryPercentageUpdate chan uint
	ramUsageUpdate          chan uint
	cpuUsageUpdate          chan uint
	wifiUpdate              chan wifiStatus
	ipUpdate                chan net.IPAddr
	transferInfoUpdate      chan transferInfo
	transferUpdater         transferUpdater
	ipUpdater               ipUpdater
}

func newStatus(tickInterval int) status {
	s := status{
		batteryPercentageUpdate: make(chan uint),
		ramUsageUpdate:          make(chan uint),
		cpuUsageUpdate:          make(chan uint),
		wifiUpdate:              make(chan wifiStatus),
		ipUpdate:                make(chan net.IPAddr),
		transferInfoUpdate:      make(chan transferInfo),
	}
	s.transferUpdater = newTransferUpdater(s.transferInfoUpdate, tickInterval)
	s.ipUpdater = newIPUpdater(s.ipUpdate, tickInterval)

	return s
}

type wifiStatus struct {
	ssid           string
	signalStrength uint
}

func (s *status) String() string {
	return fmt.Sprintf(statusBar,
		s.lastTransferInfo.String(),
		s.lastIP,
		s.lastWifi.ssid,
		s.lastWifi.signalStrength,
		s.lastBatteryPercentage,
		s.lastCPUUsage,
		s.lastRAMUsage,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

func (s *status) Run() {
	go s.transferUpdater.run()
	for {
		select {
		case bat := <-s.batteryPercentageUpdate:
			s.lastBatteryPercentage = bat
			fmt.Printf("%s\n", s.String())
		case ram := <-s.ramUsageUpdate:
			s.lastRAMUsage = ram
			fmt.Printf("%s\n", s.String())
		case cpu := <-s.cpuUsageUpdate:
			s.lastCPUUsage = cpu
			fmt.Printf("%s\n", s.String())
		case wifi := <-s.wifiUpdate:
			s.lastWifi = wifi
			fmt.Printf("%s\n", s.String())
		case ip := <-s.ipUpdate:
			s.lastIP = ip
			fmt.Printf("%s\n", s.String())
		case ti := <-s.transferInfoUpdate:
			s.lastTransferInfo = ti
			fmt.Printf("%s\n", s.String())
		}
	}
}

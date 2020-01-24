package dwmstatusgo

import (
	"fmt"
	"github.com/mdlayher/wifi"
	"log"
)

type WifiStats struct {
	client *wifi.Client
	ifi *wifi.Interface
	IfiName string
	signalPercent int
	bss *wifi.BSS
	stationInfo *wifi.StationInfo
}

func (w *WifiStats) Refresh() string {
	var err error
	w.bss, err = w.client.BSS(w.ifi)
	if err != nil {
		log.Println(err)
	}
	stationInfos, err := w.client.StationInfo(w.ifi)
	if err != nil {
		log.Println(err)
	}
	w.stationInfo = stationInfos[0]
	w.percentage()

	return fmt.Sprintf("SSID: %s %d%%(%ddBm)", w.bss.SSID, w.signalPercent, w.stationInfo.Signal)
}

func (w *WifiStats) Initialise() error {
	var err error
	w.client, err = wifi.New()
	if err != nil {
		return err
	}
	interfaces, err := w.client.Interfaces()
	if err != nil {
		return err
	}
	for _, ifi := range interfaces {
		if ifi.Name == w.IfiName {
			w.ifi = ifi
		}
	}
	return nil
}

func (w *WifiStats) percentage() {
	if w.stationInfo.Signal <= -100 {
		w.signalPercent = 0
		return
	} else if w.stationInfo.Signal >= -50 {
		w.signalPercent = 100
		return
	}
	w.signalPercent = 2 * (w.stationInfo.Signal + 100)
}

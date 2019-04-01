package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mdlayher/wifi"
)

type status string

func main() {
	//interfaces := []string{"enp4s0", "enp5s0"}
	wifiClient, err := wifi.New()
	if err != nil {
		log.Fatal(err)
	}

	interfaces, err := wifiClient.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, ifi := range interfaces {
		fmt.Printf("%s: %s, %d\n", ifi.Name, ifi.HardwareAddr, ifi.Frequency)
		si, err := wifiClient.StationInfo(ifi)
		if err == nil {
			fmt.Printf("%v\n", si)
		}
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

package query

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	iwBinary = "iw"
)

type Tools struct {
	WifiInfo *WifiInfo
}

type WifiInfo struct {
	SSID    string
	RX      int
	TX      int
	Signal  string
	Bitrate string
}

func (w WifiInfo) String() string {
	return fmt.Sprintf("SSID: %s | Signal: %s | Speed: %s", w.SSID, w.Signal, w.Bitrate)
}

func (w WifiInfo) Refresh(ifName string) error {
	cmd := exec.Command("iw", "dev", ifName, "link")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			return nil
		}
		line = strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if strings.Contains(line, "SSID") {
			w.SSID = strings.Trim(line, "SSID: ")
		} else if strings.Contains(line, "signal") {
			w.Signal = strings.Trim(line, "signal: ")
		} else if strings.Contains(line, "RX") {
			rx := strings.TrimPrefix(strings.Split(line, " bytes")[0], "RX: ")
			w.RX, err = strconv.Atoi(rx)
			if err != nil {
				return err
			}
		} else if strings.Contains(line, "TX") {
			tx := strings.TrimPrefix(strings.Split(line, " bytes")[0], "TX: ")
			w.TX, err = strconv.Atoi(tx)
			if err != nil {
				return err
			}
		} else if strings.Contains(line, "rx bitrate") {
			bitrate := strings.TrimPrefix(strings.Split(line, " VHT")[0], "rx bitrate: ")
			w.Bitrate = bitrate
		}
	}
}

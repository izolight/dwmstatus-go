package plugins

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

const (
	NMPATH    = "org.freedesktop.NetworkManager"
	SYSIFPATH = "/sys/class/net/"
)

type IPs []string

func (ips IPs) String() string {
	out := ""
	if len(ips) > 0 {
		for i, ip := range ips {
			if i != 0 {
				out += ", "
			}
			out += ip
		}
	}
	return out
}

func GetIPs(interfaces ...string) (IPs, IPs, error) {
	var ipv4s, ipv6s []string
	for _, name := range interfaces {
		netIf, err := net.InterfaceByName(name)
		if err != nil {
			return nil, nil, err
		}
		addresses, err := netIf.Addrs()
		for _, addr := range addresses {
			ipString := addr.String()
			if strings.Contains(ipString, ":") {
				if !strings.HasPrefix(ipString, "fe80:") {
					ipv6s = append(ipv6s, ipString)
				}
			} else {
				ipv4s = append(ipv4s, ipString)
			}
		}
	}
	return IPs(ipv4s), IPs(ipv6s), nil
}

type WifiInfo struct {
	SSID    string
	RX      int
	TX      int
	Signal  string
	Bitrate string
}

func (i WifiInfo) String() string {
	return fmt.Sprintf("SSID: %s | Signal: %s | Speed: %s", i.SSID, i.Signal, i.Bitrate)
}

func GetWifiInfo(wifi string) (WifiInfo, error) {
	cmd := exec.Command("iw", "dev", wifi, "link")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	info := WifiInfo{}
	if err != nil {
		return info, err
	}

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			return info, nil
		}
		line = strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if strings.Contains(line, "SSID") {
			info.SSID = strings.Trim(line, "SSID: ")
		} else if strings.Contains(line, "signal") {
			info.Signal = strings.Trim(line, "signal: ")
		} else if strings.Contains(line, "RX") {
			rx := strings.TrimPrefix(strings.Split(line, " bytes")[0], "RX: ")
			info.RX, err = strconv.Atoi(rx)
			if err != nil {
				return info, err
			}
		} else if strings.Contains(line, "TX") {
			tx := strings.TrimPrefix(strings.Split(line, " bytes")[0], "TX: ")
			info.TX, err = strconv.Atoi(tx)
			if err != nil {
				return info, err
			}
		} else if strings.Contains(line, "rx bitrate") {
			bitrate := strings.TrimPrefix(strings.Split(line, " VHT")[0], "rx bitrate: ")
			info.Bitrate = bitrate
		}
	}
}

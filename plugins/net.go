package plugins

import (
	"bytes"
	"fmt"
	"github.com/godbus/dbus"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

func GetIPs(interfaces ...string) ([]string, []string, error) {
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
	return ipv4s, ipv6s, nil
}

type WifiInfo struct {
	SSID string
	RX int
	TX int
	Signal string
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

func GetDbusPathForInterface(ifName string, conn *dbus.Conn) (dbus.ObjectPath, error) {
	var path string
	err := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager").Call("org.freedesktop.NetworkManager.GetDeviceByIpIface", 0, ifName).Store(&path)
	if err != nil {
		return "", err
	}
	return dbus.ObjectPath(path), nil
}

func GetDbusPathForAP(ifPath dbus.ObjectPath, conn *dbus.Conn) (dbus.ObjectPath, error) {
	variant, err := conn.Object("org.freedesktop.NetworkManager", ifPath).GetProperty("org.freedesktop.NetworkManager.Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return "", err
	}
	return variant.Value().(dbus.ObjectPath), nil
}

func GetSSIDFromDbus(apPath dbus.ObjectPath, conn *dbus.Conn) (string, error) {
	variant, err := conn.Object("org.freedesktop.NetworkManager", dbus.ObjectPath(apPath)).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
	if err != nil {
		return "", err
	}
	ssid := variant.Value().([]uint8)
	return string(ssid), nil
}
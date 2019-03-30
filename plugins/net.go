package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
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

func GetDbusPathForInterface(ifName string, conn *dbus.Conn) (dbus.ObjectPath, error) {
	var path string
	err := conn.Object(NMPATH, "/org/freedesktop/NetworkManager").Call(NMPATH+".GetDeviceByIpIface", 0, ifName).Store(&path)
	if err != nil {
		return "", err
	}
	return dbus.ObjectPath(path), nil
}

func GetDbusPathForAP(ifPath dbus.ObjectPath, conn *dbus.Conn) (dbus.ObjectPath, error) {
	variant, err := conn.Object(NMPATH, ifPath).GetProperty(NMPATH + ".Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return "", err
	}
	return variant.Value().(dbus.ObjectPath), nil
}

func GetSSIDFromDbus(apPath dbus.ObjectPath, conn *dbus.Conn) (string, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(apPath)).GetProperty(NMPATH + ".AccessPoint.Ssid")
	if err != nil {
		return "", err
	}
	ssid := variant.Value().([]uint8)
	return string(ssid), nil
}

func GetBitrateFromDbus(ifPath dbus.ObjectPath, conn *dbus.Conn) (uint32, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Wireless.Bitrate")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint32), nil
}

func GetTXBytesFromDbus(ifPath dbus.ObjectPath, conn *dbus.Conn) (uint64, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.TxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

func GetRXBytesFromDbus(ifPath dbus.ObjectPath, conn *dbus.Conn) (uint64, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.RxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

func GetTxBytes(ifName string) (uint64, error) {
	return getStatistics(SYSIFPATH+ifName, "tx_bytes")
}

func GetRxBytes(ifName string) (uint64, error) {
	return getStatistics(SYSIFPATH+ifName, "rx_bytes")
}

func GetRxTxBytes(ifName string) (uint64, uint64, error) {
	rx, err := getStatistics(SYSIFPATH+ifName, "rx_bytes")
	if err != nil {
		return rx, 0, err
	}
	tx, err := getStatistics(SYSIFPATH+ifName, "tx_bytes")
	return rx, tx, err
}

func getStatistics(path string, stat string) (uint64, error) {
	data, err := ioutil.ReadFile(path + "/statistics/" + stat)
	if err != nil {
		return 0, err
	}
	value := strings.TrimSuffix(string(data), "\n")
	tx, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return tx, nil
}

package query

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// SysPath is the path where the /sys filesystem is located (mostly used to change it for tests)
var SysPath = "/sys"

const (
	ifPath      = "/class/net/"
	batteryPath = "/class/power_supply/"
)

type SysFs struct {
	BatteryInfo   map[string]*BatteryInfo
	InterfaceInfo map[string]*InterfaceInfo
}

type BatteryInfo struct {
	MaxBatteryCapacity     uint64
	CurrentBatteryCapacity uint64
}

// ReadUint64 returns a uint64 read from a file (in /sys)
func ReadUint64(path string) (uint64, error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	fileValue := strings.TrimSuffix(string(fileData), "\n")
	value, err := strconv.ParseUint(fileValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (s *SysFs) RefreshCurrentBatteryCapacity(batteries ...string) error {
	for _, b := range batteries {
		cap, err := ReadUint64(SysPath + batteryPath + b + "/energy_now")
		if err != nil {
			return err
		}
		s.BatteryInfo[b].CurrentBatteryCapacity = cap
	}
	return nil

}

func (s *SysFs) RefreshMaxBatteryCapacity(batteries ...string) error {
	for _, b := range batteries {
		cap, err := ReadUint64(SysPath + batteryPath + b + "/energy_full")
		if err != nil {
			return err
		}
		s.BatteryInfo[b].MaxBatteryCapacity = cap
	}
	return nil
}

// RefreshTxBytes returns the transmitted bytes for an interface
func (s *SysFs) RefreshTxBytes(ifName string) error {
	txBytes, err := ReadUint64(SysPath + ifPath + ifName + "/statistics/tx_bytes")
	if err != nil {
		return err
	}
	s.InterfaceInfo[ifName].TxBytes = txBytes
	return nil
}

// RefreshRxBytes returns the received bytes for an interface
func (s *SysFs) RefreshRxBytes(ifName string) error {
	rxBytes, err := ReadUint64(SysPath + ifPath + ifName + "/statistics/rx_bytes")
	if err != nil {
		return err
	}
	s.InterfaceInfo[ifName].RxBytes = rxBytes
	return nil
}

// RefreshRxTxBytes combines the Rx and Tx Bytes functions
func (s *SysFs) RefreshRxTxBytes(ifName string) error {
	err := s.RefreshRxBytes(ifName)
	if err != nil {
		return err
	}
	err = s.RefreshTxBytes(ifName)
	if err != nil {
		return err
	}
	return nil
}

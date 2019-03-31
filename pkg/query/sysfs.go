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

// RefreshAll refreshes data from all registered plugins
func (s *SysFs) RefreshAll() error {
	err := s.RefreshCurrentBatteryCapacity()
	if err != nil {
		return err
	}
	err = s.RefreshMaxBatteryCapacity()
	if err != nil {
		return err
	}
	err = s.RefreshTxBytes()
	if err != nil {
		return err
	}
	err = s.RefreshRxBytes()
	if err != nil {
		return err
	}
	return nil
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

// RefreshCurrentBatteryCapacity refreshes the current battery capacity from sysfs
func (s *SysFs) RefreshCurrentBatteryCapacity() error {
	for b := range s.BatteryInfo {
		cap, err := ReadUint64(SysPath + batteryPath + b + "/energy_now")
		if err != nil {
			return err
		}
		s.BatteryInfo[b].CurrentBatteryCapacity = cap
	}
	return nil

}

// RefreshMaxBatteryCapacity refreshes the max battery capacity from sysfs
func (s *SysFs) RefreshMaxBatteryCapacity() error {
	for b := range s.BatteryInfo {
		cap, err := ReadUint64(SysPath + batteryPath + b + "/energy_full")
		if err != nil {
			return err
		}
		s.BatteryInfo[b].MaxBatteryCapacity = cap
	}
	return nil
}

// RefreshTxBytes refreshes the transmitted bytes from sysfs
func (s *SysFs) RefreshTxBytes() error {
	for i := range s.InterfaceInfo {
		txBytes, err := ReadUint64(SysPath + ifPath + i + "/statistics/tx_bytes")
		if err != nil {
			return err
		}
		s.InterfaceInfo[i].TxBytes = txBytes
	}
	return nil
}

// RefreshRxBytes returns the received bytes for an interface
func (s *SysFs) RefreshRxBytes() error {
	for i := range s.InterfaceInfo {
		rxBytes, err := ReadUint64(SysPath + ifPath + i + "/statistics/rx_bytes")
		if err != nil {
			return err
		}
		s.InterfaceInfo[i].RxBytes = rxBytes
	}
	return nil
}

// RefreshRxTxBytes combines the Rx and Tx Bytes functions
func (s *SysFs) RefreshRxTxBytes() error {
	err := s.RefreshRxBytes()
	if err != nil {
		return err
	}
	err = s.RefreshTxBytes()
	if err != nil {
		return err
	}
	return nil
}

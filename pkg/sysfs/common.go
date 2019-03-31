package sysfs

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

// Uint64 returns a uint64 read from a file (in /sys)
func Uint64(path string) (uint64, error) {
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

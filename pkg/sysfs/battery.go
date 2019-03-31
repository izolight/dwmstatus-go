package sysfs

import (
	"fmt"
	"time"
)

func CurrentBatteryCapacity(batteries ...string) (uint64, error) {
	var capacity uint64
	for _, b := range batteries {
		cap, err := Uint64(SysPath + batteryPath + b + "/energy_now")
		if err != nil {
			return capacity, err
		}
		capacity += cap
	}
	return capacity, nil

}

func MaxBatteryCapacity(batteries ...string) (uint64, error) {
	var capacity uint64
	for _, b := range batteries {
		cap, err := Uint64(SysPath + batteryPath + b + "/energy_full")
		if err != nil {
			return capacity, err
		}
		capacity += cap
	}
	return capacity, nil
}

func CalculateRemainingTime(current uint64, previous uint64, duration time.Duration) uint64 {
	difference := current - previous
	return (current / difference) * uint64(duration)
}

func CalculateBatteryPercentage(current uint64, max uint64) (uint8, error) {
	if max < current {
		return 0, fmt.Errorf("Max: %d has to be greater or equal than current: %d", max, current)
	}
	return uint8(current / max * 100), nil
}

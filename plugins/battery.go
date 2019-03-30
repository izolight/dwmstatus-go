package plugins

import (
	"fmt"
	"time"
)

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

package dwmstatusgo

import (
	"fmt"
	"github.com/prometheus/procfs/sysfs"
	"log"
	"time"
)

type BatteryStats struct {
	Batteries []string
	energyNow int64
	energyFull int64
	energyFullDesign int64
	energyNowPrevious int64
	remaining int64
	nextMeasurement time.Time
}

func (b *BatteryStats)calculateRemainingTime() int64 {
	remainingEnergy := b.energyFull - b.energyNow
	charged := b.energyNow - b.energyNowPrevious
	if charged == 0 {
		return 0
	} else if charged > 0 {
		minutes := float64(remainingEnergy) / float64(charged) / 60 * 5
		return int64(minutes)
	} else {
		minutes := float64(b.energyNow) / float64(charged) / 60 * 5
		return int64(minutes) * (-1)
	}
}

func (b *BatteryStats) Refresh() string {
	psc, err := sysfs.NewPowerSupplyClass()
	if err != nil {
		log.Println(err)
	}
	for _, v := range b.Batteries {
		bat, ok := psc[v]
		if !ok {
			continue
		}
		b.energyNow += *bat.EnergyNow
		b.energyFull += *bat.EnergyFull
		b.energyFullDesign += *bat.EnergyFullDesign
	}
	measure := time.Now().Truncate(time.Second)
	if measure.UnixNano() == b.nextMeasurement.UnixNano() {
		b.nextMeasurement = b.nextMeasurement.Add(5 * time.Second)
		b.remaining = b.calculateRemainingTime()
		b.energyNowPrevious = b.energyNow
	}
	return fmt.Sprintf("BAT: %.1f%% %dMin", (float64(b.energyNow)/float64(b.energyFull))*100, b.remaining)
}
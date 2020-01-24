package dwmstatus

import (
	"fmt"

	"github.com/prometheus/procfs"
)

type CPUInfo struct {
	fs           procfs.FS
	lastCPUStat  procfs.CPUStat
	lastCPUUsage float64
}

func NewCPUInfo() (*CPUInfo, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return nil, fmt.Errorf("couldn't get new procfs: %w", err)
	}
	stat, err := fs.NewStat()
	if err != nil {
		return nil, fmt.Errorf("couldn't stat: %w", err)
	}
	return &CPUInfo{
		fs:           fs,
		lastCPUStat:  stat.CPUTotal,
		lastCPUUsage: 0,
	}, nil
}

func (c *CPUInfo) Refresh() string {
	stat, err := c.fs.NewStat()
	if err != nil {
		return fmt.Sprintf("%.3f%%", c.lastCPUUsage)
	}
	lastCPUStat := c.lastCPUStat
	c.lastCPUStat = stat.CPUTotal
	c.lastCPUUsage = usageCPUStat(subtract(c.lastCPUStat, lastCPUStat))

	return fmt.Sprintf("%.3f%%", c.lastCPUUsage)
}

func sumCPUStat(c procfs.CPUStat) float64 {
	return c.User + c.Nice + c.System + c.Idle + c.Iowait + c.IRQ + c.SoftIRQ + c.Steal + c.Guest + c.GuestNice
}

func usageCPUStat(c procfs.CPUStat) float64 {
	return (1 - (c.Idle / sumCPUStat(c))) * 100
}

func subtract(a, b procfs.CPUStat) procfs.CPUStat {
	return procfs.CPUStat{
		User:      a.User - b.User,
		Nice:      a.Nice - b.Nice,
		System:    a.System - b.System,
		Idle:      a.Idle - b.Idle,
		Iowait:    a.Iowait - b.Iowait,
		IRQ:       a.IRQ - b.IRQ,
		SoftIRQ:   a.SoftIRQ - b.SoftIRQ,
		Steal:     a.Steal - b.Steal,
		Guest:     a.Guest - b.Guest,
		GuestNice: a.GuestNice - b.GuestNice,
	}
}

package dwmstatus

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
)

type MemoryInfo struct {
	virtualMem *mem.VirtualMemoryStat
}

func NewMemoryInfo() *MemoryInfo {
	m := &MemoryInfo{}
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	m.virtualMem = vm
	return m
}

func (m *MemoryInfo) Refresh() string {
	info := m.virtualMem

	return fmt.Sprintf("RAM: %s", humanize.Bytes(info.Used))
}

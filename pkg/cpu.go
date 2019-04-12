package dwmstatusgo

import (
	"fmt"
	"github.com/prometheus/procfs"
	"log"
)

type CPUInfo struct {}

func (c *CPUInfo) Refresh() string {
	fs, err := procfs.NewFS("/")
	if err != nil {
		log.Println(err)
	}
	stat, err := fs.NewStat()
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("CPU: %.3f", stat.CPUTotal.System)
}
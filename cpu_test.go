package dwmstatus

import (
	"fmt"
	"testing"
	"time"
)

func TestRefreshCpuInfo(t *testing.T) {
	c, _ := NewCPUInfo()
	for i := 0; i < 10; i++ {
		fmt.Println(c.Refresh())
		time.Sleep(1 * time.Second)
	}
}

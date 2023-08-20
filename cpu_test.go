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

func TestParseCPUStat(t *testing.T) {
	p := NewCPUStat(WithPath("testdata/procstat"))
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p)
}
package dwmstatus

import "testing"

func TestStatusString(t *testing.T) {
	s := newStatus(1)
	s.Run()
}

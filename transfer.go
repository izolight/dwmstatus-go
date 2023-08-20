package dwmstatus

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/procfs"
)

type transferInfo struct {
	received    uint64
	transferred uint64
}

func (t *transferInfo) String() string {
	return fmt.Sprintf("U: %d D: %d", t.transferred, t.received)
}

type transferUpdater struct {
	lastTransfer       transferInfo
	transferInfoUpdate chan transferInfo
	ticker             *time.Ticker
}

func newTransferUpdater(transferUpdate chan transferInfo, tickInterval int) transferUpdater {
	t := transferUpdater{
		transferInfoUpdate: transferUpdate,
		ticker:             time.NewTicker(time.Duration(tickInterval) * time.Second),
	}
	nd, err := procfs.NewNetDev()
	if err != nil {
		return t
	}
	total := nd.Total()
	t.lastTransfer.received = total.RxBytes
	t.lastTransfer.transferred = total.TxBytes

	return t
}

func (t *transferUpdater) run() {
	for {
		select {
		case <-t.ticker.C:
			nd, err := procfs.NewNetDev()
			if err != nil {
				return
			}
			total := nd.Total()
			t.transferInfoUpdate <- transferInfo{
				received:    total.RxBytes - t.lastTransfer.received,
				transferred: total.TxBytes - t.lastTransfer.transferred,
			}

			t.lastTransfer.received = total.RxBytes
			t.lastTransfer.transferred = total.TxBytes
		}
	}
}

func updateTransfer() {
	var file *os.File
	var err error

	for {
		file, err = os.Open("/proc/net/dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
		var void = 0
		var dev, rx, tx = "", 0, 0
		scanner := bufio.NewScanner(file)
		// run scan 2 times to skip the header
		scanner.Scan() 
		scanner.Scan()
		for scanner.Scan() {
			_, err = fmt.Sscanf(scanner.Text(), "%s %d %d %d %d %d %d %d %d %d",
				&dev, &rx, &void, &void, &void, &void, &void, &void, &void, &tx)
			fmt.Println(dev, rx, tx)
		}
		time.Sleep(time.Second)
	}
}

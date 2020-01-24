package dwmstatus

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
	"github.com/prometheus/procfs"
)

type TransferInfo struct {
	prevRx uint64
	prevTx uint64
}

func (t *TransferInfo) Refresh() string {
	nd, err := procfs.NewNetDev()
	if err != nil {
		log.Println(err)
	}
	total := nd.Total()
	rx, tx := total.RxBytes, total.TxBytes

	out := fmt.Sprintf("U:%s/s D:%s/s", humanize.Bytes(tx-t.prevTx), humanize.Bytes(rx-t.prevRx))
	t.prevRx, t.prevTx = rx, tx
	return out
}

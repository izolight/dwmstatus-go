package sysfs

// TxBytes returns the transmitted bytes for an interface
func TxBytes(ifName string) (uint64, error) {
	return Uint64(SysPath + ifPath + ifName + "/statistics/tx_bytes")
}

// RxBytes returns the received bytes for an interface
func RxBytes(ifName string) (uint64, error) {
	return Uint64(SysPath + ifPath + ifName + "/statistics/rx_bytes")
}

// RxTxBytes combines the Rx and Tx Bytes functions
func RxTxBytes(ifName string) (uint64, uint64, error) {
	rx, err := Uint64(SysPath + ifPath + ifName + "/statistics/rx_bytes")
	if err != nil {
		return rx, 0, err
	}
	tx, err := Uint64(SysPath + ifPath + ifName + "/statistics/tx_bytes")
	return rx, tx, err
}

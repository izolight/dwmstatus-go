package dbus

import "github.com/godbus/dbus"

func SSID(apPath dbus.ObjectPath) (string, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(apPath)).GetProperty(NMPATH + ".AccessPoint.Ssid")
	if err != nil {
		return "", err
	}
	ssid := variant.Value().([]uint8)
	return string(ssid), nil
}

func WifiLinkSpeed(ifPath dbus.ObjectPath) (uint32, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Wireless.Bitrate")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint32), nil
}

func TxBytes(ifPath dbus.ObjectPath) (uint64, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.TxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

func RxBytes(ifPath dbus.ObjectPath) (uint64, error) {
	variant, err := conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.RxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

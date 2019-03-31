package query

import "github.com/godbus/dbus"

type Dbus struct {
	conn *dbus.Conn
}

const (
	NMPATH = "org.freedesktop.NetworkManager"
)

func (d *Dbus) Connect() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

func (d *Dbus) PathForInterface(ifName string) (dbus.ObjectPath, error) {
	path := ""
	err := d.conn.Object(NMPATH, "/"+NMPATH).Call(NMPATH+".GetDeviceByIpIface", 0, ifName).Store(&path)
	if err != nil {
		return "", err
	}
	return dbus.ObjectPath(path), nil
}

func (d *Dbus) PathForAccessPoint(ifPath dbus.ObjectPath) (dbus.ObjectPath, error) {
	variant, err := d.conn.Object(NMPATH, ifPath).GetProperty(NMPATH + ".Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return "", err
	}
	return variant.Value().(dbus.ObjectPath), nil
}

func (d *Dbus) SSID(apPath dbus.ObjectPath) (string, error) {
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(apPath)).GetProperty(NMPATH + ".AccessPoint.Ssid")
	if err != nil {
		return "", err
	}
	ssid := variant.Value().([]uint8)
	return string(ssid), nil
}

func (d *Dbus) WifiLinkSpeed(ifPath dbus.ObjectPath) (uint32, error) {
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Wireless.Bitrate")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint32), nil
}

func (d *Dbus) TxBytes(ifPath dbus.ObjectPath) (uint64, error) {
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.TxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

func (d *Dbus) RxBytes(ifPath dbus.ObjectPath) (uint64, error) {
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(ifPath)).GetProperty(NMPATH + ".Device.Statistics.RxBytes")
	if err != nil {
		return 0, err
	}
	return variant.Value().(uint64), nil
}

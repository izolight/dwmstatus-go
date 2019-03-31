package dbus

import (
	"log"

	"github.com/godbus/dbus"
)

var conn dbus.Conn

const (
	NMPATH = "org.freedesktop.NetworkManager"
)

func init() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Couldn't connect to Systembus (%v): %s", conn, err)
	}
}

func PathForInterface(ifName string) (dbus.ObjectPath, error) {
	path := ""
	err := conn.Object(NMPATH, "/"+NMPATH).Call(NMPATH+".GetDeviceByIpIface", 0, ifName).Store(&path)
	if err != nil {
		return "", err
	}
	return dbus.ObjectPath(path), nil
}

func PathForAccessPoint(ifPath dbus.ObjectPath) (dbus.ObjectPath, error) {
	variant, err := conn.Object(NMPATH, ifPath).GetProperty(NMPATH + ".Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return "", err
	}
	return variant.Value().(dbus.ObjectPath), nil
}

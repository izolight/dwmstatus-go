package query

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus"
)

type Dbus struct {
	conn          *dbus.Conn
	InterfaceInfo map[string]*InterfaceInfo
}

type InterfaceInfo struct {
	SSID      string
	LinkSpeed uint32
	TxBytes   uint64
	RxBytes   uint64
	ifPath    dbus.ObjectPath
	apPath    dbus.ObjectPath
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

func (d *Dbus) InterfacePath(ifName string) error {
	path := ""
	err := d.conn.Object(NMPATH, "/"+NMPATH).Call(NMPATH+".GetDeviceByIpIface", 0, ifName).Store(&path)
	if err != nil {
		return err
	}
	_, ok := d.InterfaceInfo[ifName]
	if !ok {
		d.InterfaceInfo[ifName] = &InterfaceInfo{}
	}
	d.InterfaceInfo[ifName].ifPath = dbus.ObjectPath(path)
	return nil
}

func (d *Dbus) AccessPointPath(ifName string) error {
	_, ok := d.InterfaceInfo[ifName]
	if !ok || ok && d.InterfaceInfo[ifName].ifPath == "" {
		return errors.New("You need to set the Interface Path first")
	}
	variant, err := d.conn.Object(NMPATH, d.InterfaceInfo[ifName].ifPath).GetProperty(NMPATH + ".Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return err
	}
	apPath, ok := variant.Value().(dbus.ObjectPath)
	if !ok {
		return fmt.Errorf("Couldn't get valid ObjectPath for ap: %s", variant.Value())
	}
	d.InterfaceInfo[ifName].apPath = apPath
	return nil
}

func (d *Dbus) RefreshSSID(ifName string) error {
	_, ok := d.InterfaceInfo[ifName]
	if !ok || ok && d.InterfaceInfo[ifName].apPath == "" {
		return errors.New("You need to set the Access Point Path first")
	}
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(d.InterfaceInfo[ifName].apPath)).GetProperty(NMPATH + ".AccessPoint.Ssid")
	if err != nil {
		return err
	}
	ssid, ok := variant.Value().([]uint8)
	if !ok {
		return fmt.Errorf("Couldn't get valid SSID: %s", variant.Value())
	}
	d.InterfaceInfo[ifName].SSID = string(ssid)
	return nil
}

func (d *Dbus) WifiLinkSpeed(ifName string) error {
	_, ok := d.InterfaceInfo[ifName]
	if !ok || ok && d.InterfaceInfo[ifName].ifPath == "" {
		return errors.New("You need to set the InterfacePath Path first")
	}
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(d.InterfaceInfo[ifName].ifPath)).GetProperty(NMPATH + ".Device.Wireless.Bitrate")
	if err != nil {
		return err
	}
	linkSpeed, ok := variant.Value().(uint32)
	if !ok {
		return fmt.Errorf("Couldn't get valid LinkSpeed: %s", variant.Value())
	}
	d.InterfaceInfo[ifName].LinkSpeed = linkSpeed
	return nil
}

func (d *Dbus) TxBytes(ifName string) error {
	_, ok := d.InterfaceInfo[ifName]
	if !ok || ok && d.InterfaceInfo[ifName].ifPath == "" {
		return errors.New("You need to set the InterfacePath Path first")
	}
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(d.InterfaceInfo[ifName].ifPath)).GetProperty(NMPATH + ".Device.Statistics.TxBytes")
	if err != nil {
		return err
	}
	txBytes, ok := variant.Value().(uint64)
	if !ok {
		return fmt.Errorf("Couldn't get valid TxBytes: %s", variant.Value())
	}
	d.InterfaceInfo[ifName].TxBytes = txBytes
	return nil
}

func (d *Dbus) RxBytes(ifName string) error {
	_, ok := d.InterfaceInfo[ifName]
	if !ok || ok && d.InterfaceInfo[ifName].ifPath == "" {
		return errors.New("You need to set the InterfacePath Path first")
	}
	variant, err := d.conn.Object(NMPATH, dbus.ObjectPath(d.InterfaceInfo[ifName].ifPath)).GetProperty(NMPATH + ".Device.Statistics.RxBytes")
	if err != nil {
		return err
	}
	rxBytes, ok := variant.Value().(uint64)
	if !ok {
		return fmt.Errorf("Couldn't get valid RxBytes: %s", variant.Value())
	}
	d.InterfaceInfo[ifName].RxBytes = rxBytes
	return nil
}

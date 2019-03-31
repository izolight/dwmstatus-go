package query_test

import (
	"testing"

	"github.com/izolight/dwmstatus-go/pkg/query"
)

func init() {
	query.SysPath = "./testdata"
}

func TestUint64(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"Not a uint64", args{path: "./testdata/class/net/fake0/statistics/notUint64"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := query.ReadUint64(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSysFs_RefreshCurrentBatteryCapacity(t *testing.T) {
	type args struct {
		battery string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"ExistingBattery", args{"BAT0"}, 120080, false},
		{"NonExistingBattery", args{"BAT1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &query.SysFs{
				BatteryInfo: map[string]*query.BatteryInfo{
					tt.args.battery: &query.BatteryInfo{},
				},
			}
			if err := s.RefreshCurrentBatteryCapacity(tt.args.battery); (err != nil) != tt.wantErr {
				t.Errorf("SysFs.RefreshCurrentBatteryCapacity() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.BatteryInfo[tt.args.battery].CurrentBatteryCapacity
			if got != tt.want {
				t.Errorf("SysFs.RefreshCurrentBatteryCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSysFs_RefreshMaxBatteryCapacity(t *testing.T) {
	type args struct {
		battery string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"ExistingBattery", args{"BAT0"}, 500600, false},
		{"NonExistingBattery", args{"BAT1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &query.SysFs{
				BatteryInfo: map[string]*query.BatteryInfo{
					tt.args.battery: &query.BatteryInfo{},
				},
			}
			if err := s.RefreshMaxBatteryCapacity(tt.args.battery); (err != nil) != tt.wantErr {
				t.Errorf("SysFs.RefreshMaxBatteryCapacity() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.BatteryInfo[tt.args.battery].MaxBatteryCapacity
			if got != tt.want {
				t.Errorf("SysFs.RefreshMaxBatteryCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSysFs_RefreshTxBytes(t *testing.T) {
	type args struct {
		ifName string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"fakeInterface", args{ifName: "fake0"}, 1234, false},
		{"fakeInterfaceNotFound", args{ifName: "fake1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &query.SysFs{
				InterfaceInfo: map[string]*query.InterfaceInfo{
					tt.args.ifName: &query.InterfaceInfo{},
				},
			}
			if err := s.RefreshTxBytes(tt.args.ifName); (err != nil) != tt.wantErr {
				t.Errorf("SysFs.RefreshTxBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.InterfaceInfo[tt.args.ifName].TxBytes
			if got != tt.want {
				t.Errorf("SysFs.RefreshTxBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSysFs_RefreshRxBytes(t *testing.T) {
	type args struct {
		ifName string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"fakeInterface", args{ifName: "fake0"}, 9876, false},
		{"fakeInterfaceNotFound", args{ifName: "fake1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &query.SysFs{
				InterfaceInfo: map[string]*query.InterfaceInfo{
					tt.args.ifName: &query.InterfaceInfo{},
				},
			}
			if err := s.RefreshRxBytes(tt.args.ifName); (err != nil) != tt.wantErr {
				t.Errorf("SysFs.RefreshRxBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.InterfaceInfo[tt.args.ifName].RxBytes
			if got != tt.want {
				t.Errorf("SysFs.RefreshRxBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSysFs_RefreshRxTxBytes(t *testing.T) {
	type args struct {
		ifName string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		want1   uint64
		wantErr bool
	}{
		{"fakeInterface", args{ifName: "fake0"}, 9876, 1234, false},
		{"fakeInterface", args{ifName: "fake1"}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &query.SysFs{
				InterfaceInfo: map[string]*query.InterfaceInfo{
					tt.args.ifName: &query.InterfaceInfo{},
				},
			}
			if err := s.RefreshRxTxBytes(tt.args.ifName); (err != nil) != tt.wantErr {
				t.Errorf("SysFs.RefreshRxTxBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.InterfaceInfo[tt.args.ifName].RxBytes
			if got != tt.want {
				t.Errorf("SysFs.RefreshRxTxBytes() = %v, want %v", got, tt.want)
			}
			got1 := s.InterfaceInfo[tt.args.ifName].TxBytes
			if got1 != tt.want1 {
				t.Errorf("SysFs.RefreshRxTxBytes() = %v, want %v", got1, tt.want1)
			}
		})
	}
}

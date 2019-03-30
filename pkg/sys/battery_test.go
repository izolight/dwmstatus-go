package sys_test

import (
	"testing"

	"github.com/izolight/dwmstatus-go/pkg/sys"
)

func TestCurrentBatterCapacity(t *testing.T) {
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
			got, err := sys.CurrentBatteryCapacity(tt.args.battery)
			if (err != nil) != tt.wantErr {
				t.Errorf("CurrentBatteryCapacity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CurrentBatteryCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxBatteryCapacity(t *testing.T) {
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
			got, err := sys.MaxBatteryCapacity(tt.args.battery)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxBatteryCapacity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MaxBatteryCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

package sys_test

import (
	"testing"

	"github.com/izolight/dwmstatus-go/pkg/sys"
)

func TestTxBytes(t *testing.T) {
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
			got, err := sys.TxBytes(tt.args.ifName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TxBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRxBytes(t *testing.T) {
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
			got, err := sys.RxBytes(tt.args.ifName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RxBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RxBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRxTxBytes(t *testing.T) {
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
			got, got1, err := sys.RxTxBytes(tt.args.ifName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RxTxBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RxTxBytes() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RxTxBytes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

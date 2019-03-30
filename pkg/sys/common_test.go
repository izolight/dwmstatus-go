package sys_test

import (
	"testing"

	"github.com/izolight/dwmstatus-go/pkg/sys"
)

func init() {
	sys.SysPath = "./testdata"
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
			got, err := sys.Uint64(tt.args.path)
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

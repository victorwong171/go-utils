package logger

import (
	"reflect"
	"testing"

	"github.com/victorwong171/fatwang-go-utils/utils"
)

func (n *noop) WithField(key string, value interface{}) utils.Logger {
	return n
}

func (n *noop) Debugf(format string, args ...interface{}) {}

func (n *noop) Infof(format string, args ...interface{}) {}

func (n *noop) Warnf(format string, args ...interface{}) {}

func (n *noop) Errorf(format string, args ...interface{}) {}

func (n *noop) Fatalf(format string, args ...interface{}) {}

func (n *noop) Panicf(format string, args ...interface{}) {}

func (n *noop) Debug(args ...interface{}) {}

func (n *noop) Info(args ...interface{}) {}

func (n *noop) Warn(args ...interface{}) {}

func (n *noop) Error(args ...interface{}) {}

func (n *noop) Fatal(args ...interface{}) {}

func (n *noop) Panic(args ...interface{}) {}

func (n *noop) Flush() error {
	return nil
}
func TestNewNoop(t *testing.T) {
	tests := []struct {
		name string
		want *noop
	}{
		{
			name: "NewNoop",
			want: &noop{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_noop_Info(t *testing.T) {
	type fields struct {
	}
	type args struct {
		in0 string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "all is ok",
			fields: fields{},
			args:   args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &noop{}
			n.Info(tt.args.in0)
		})
	}
}

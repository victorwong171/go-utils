package bitmap

import (
	"math"
	"reflect"
	"testing"
)

var (
	bm = NewBitMap(math.MaxInt16 >> 1)
)

func BenchmarkBitMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bm.Set(uint(i % (math.MaxInt16 >> 1)))
	}
}

func BenchmarkBitMap_Check(b *testing.B) {
	bm.Set(1024)
	for i := 0; i < b.N; i++ {
		bm.Check(1024)
	}
}

func BenchmarkBitMap_Del(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bm.Del(uint(i % (math.MaxInt16 >> 1)))
	}
}

var (
	defaultMax uint = 8192
	size            = (defaultMax + 7) >> 3
)

func TestNewBitMap(t *testing.T) {
	type args struct {
		max uint
	}
	tests := []struct {
		name string
		args args
		want *bitMap
	}{
		{
			name: "all is ok",
			args: args{
				max: defaultMax,
			},
			want: &bitMap{vmax: defaultMax, graph: make([]byte, size, size)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitMap(tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bitMap_Check(t *testing.T) {
	type fields struct {
		vmax  uint
		graph []byte
	}
	type args struct {
		num uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "all is ok",
			fields: fields{
				vmax:  defaultMax,
				graph: make([]byte, size, size),
			},
			args: args{
				num: 1024,
			},
			want: false,
		},
		{
			name: "oversize",
			fields: fields{
				vmax:  defaultMax,
				graph: make([]byte, size, size),
			},
			args: args{
				num: 8193,
			},
			want: false,
		},
		{
			name: "check bit that is set",
			fields: fields{
				vmax:  defaultMax,
				graph: func() []byte { g := make([]byte, size, size); g[1] = 4; return g }(),
			},
			args: args{
				num: 10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bitMap{
				vmax:  tt.fields.vmax,
				graph: tt.fields.graph,
			}
			if got := b.Check(tt.args.num); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bitMap_Del(t *testing.T) {
	type fields struct {
		vmax  uint
		graph []byte
	}
	type args struct {
		num uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "all is ok",
			fields: fields{
				vmax:  defaultMax,
				graph: make([]byte, size, size),
			},
			args: args{
				num: 1024,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bitMap{
				vmax:  tt.fields.vmax,
				graph: tt.fields.graph,
			}
			b.Del(tt.args.num)
		})
	}
}

func Test_bitMap_Set(t *testing.T) {
	type fields struct {
		vmax  uint
		graph []byte
	}
	type args struct {
		num uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "all is ok",
			fields: fields{
				vmax:  defaultMax,
				graph: make([]byte, size, size),
			},
			args: args{
				num: 1024,
			},
		},
		{
			name: "oversize",
			fields: fields{
				vmax:  defaultMax,
				graph: make([]byte, size, size),
			},
			args: args{
				num: 9217,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bitMap{
				vmax:  tt.fields.vmax,
				graph: tt.fields.graph,
			}
			b.Set(tt.args.num)
		})
	}
}

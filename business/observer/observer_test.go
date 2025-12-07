package observer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/victorwong171/fatwang-go-utils/utils"
)

func TestNewEvent(t *testing.T) {
	type args[P any] struct {
		observers []Cfg[int]
	}
	cfgList := []Cfg[int]{
		{
			IsAsync: true,
			Name:    "async failed",
			Action: func(ctx context.Context, params int) error {
				return errors.New("")
			},
		},
		{
			IsAsync: false,
			Name:    "sync ok",
			Action: func(ctx context.Context, params int) error {
				return nil
			},
		},
	}

	type testCase[P any] struct {
		name string
		args args[P]
		want Event[int]
	}
	tests := []testCase[int]{
		{
			name: "all is ok",
			args: args[int]{
				observers: cfgList,
			},
			want: &event[int]{
				observerList: cfgList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEvent(utils.MustNewDevelopment(), tt.args.observers...)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_event_Emit(t *testing.T) {
	type args[P any] struct {
		ctx    context.Context
		params P
	}
	type testCase[P any] struct {
		name    string
		e       event[P]
		args    args[P]
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "all is ok",
			e: event[int]{
				logger: utils.MustNewDevelopment(),
				observerList: []Cfg[int]{
					{
						IsAsync: true,
						Name:    "async failed",
						Action: func(ctx context.Context, params int) error {
							return errors.New("")
						},
					},
					{
						IsAsync: false,
						Name:    "sync ok",
						Action: func(ctx context.Context, params int) error {
							time.Sleep(1 * time.Second)
							return nil
						},
					},
				},
			},
			args: args[int]{
				ctx:    context.Background(),
				params: 0,
			},
			wantErr: false,
		},
		{
			name: "sync failed",
			e: event[int]{
				logger: utils.MustNewDevelopment(),
				observerList: []Cfg[int]{
					{
						IsAsync: true,
						Name:    "async failed",
						Action: func(ctx context.Context, params int) error {
							return errors.New("")
						},
					},
					{
						IsAsync: false,
						Name:    "sync ok",
						Action: func(ctx context.Context, params int) error {

							return errors.New("")
						},
					},
				},
			},
			args: args[int]{
				ctx:    context.Background(),
				params: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Emit(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Emit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_event_Register(t *testing.T) {
	type args[P any] struct {
		observers []Cfg[int]
	}
	type testCase[P any] struct {
		name string
		e    event[P]
		args args[P]
	}
	tests := []testCase[int]{
		{
			name: "all is ok",
			e: event[int]{
				observerList: []Cfg[int]{{}},
			},
			args: args[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Register(tt.args.observers...)
		})
	}
}

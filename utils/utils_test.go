package utils

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go"
)

func TestGetUuid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ok",
			want: GetUuid(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUuid(); len(got) != len(tt.want) {
				t.Errorf("GetUuid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernaryOperator(t *testing.T) {
	type args[T any] struct {
		expr     bool
		trueVal  T
		falseVal T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[bool]{
		{
			name: "true",
			args: args[bool]{
				expr:     1 == 1,
				trueVal:  true,
				falseVal: false,
			},
			want: true,
		},
		{
			name: "false",
			args: args[bool]{
				expr:     1 != 1,
				trueVal:  true,
				falseVal: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TernaryOperator(tt.args.expr, tt.args.trueVal, tt.args.falseVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TernaryOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextCopy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "all is ok",
			args: args{
				ctx: ctx,
			},
			want: zipkin.NewContext(ctx, zipkin.SpanOrNoopFromContext(ctx)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContextCopy(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContextCopy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomIntList(t *testing.T) {
	type args struct {
		size  int
		edge1 int
		edge2 int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all is ok",
			args: args{
				size:  3,
				edge1: 3,
				edge2: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRandomIntList(tt.args.size, tt.args.edge1, tt.args.edge2)
			if len(got) != tt.args.size {
				t.Errorf("GenerateRandomIntList() length = %v, want %v", len(got), tt.args.size)
			}
			for _, v := range got {
				if v < tt.args.edge2 || v > tt.args.edge1 {
					t.Errorf("GenerateRandomIntList() value = %v, want between %v and %v", v, tt.args.edge2, tt.args.edge1)
				}
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args[T NUMBER] struct {
		m T
		n T
	}
	type testCase[T NUMBER] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "all is ok",
			args: args[int]{
				m: 1,
				n: 0,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.m, tt.args.n); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args[T NUMBER] struct {
		m T
		n T
	}
	type testCase[T NUMBER] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "all is ok",
			args: args[int]{
				m: 1,
				n: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.m, tt.args.n); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentLimit(t *testing.T) {
	// 测试正常情况
	t.Run("normal case", func(t *testing.T) {
		var results []int
		var mu sync.Mutex
		
		workList := []func() error{
			func() error {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				results = append(results, 1)
				mu.Unlock()
				return nil
			},
			func() error {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				results = append(results, 2)
				mu.Unlock()
				return nil
			},
		}
		
		err := CurrentLimit(2, workList, false)
		if err != nil {
			t.Errorf("CurrentLimit() error = %v, want nil", err)
		}
		
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})
	
	// 测试限制并发数
	t.Run("limit concurrency", func(t *testing.T) {
		var results []int
		var mu sync.Mutex
		start := time.Now()
		
		workList := []func() error{
			func() error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				results = append(results, 1)
				mu.Unlock()
				return nil
			},
			func() error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				results = append(results, 2)
				mu.Unlock()
				return nil
			},
			func() error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				results = append(results, 3)
				mu.Unlock()
				return nil
			},
		}
		
		// 限制并发数为2
		err := CurrentLimit(2, workList, false)
		duration := time.Since(start)
		if err != nil {
			t.Errorf("CurrentLimit() error = %v, want nil", err)
		}
		
		// 由于并发限制，执行时间应该大于100ms（3个任务，每批2个）
		if duration < 100*time.Millisecond {
			t.Errorf("Expected duration > 100ms, got %v", duration)
		}
		
		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}
	})
	
	// 测试错误处理 enableError=false
	t.Run("error handling with enableError=false", func(t *testing.T) {
		workList := []func() error{
			func() error {
				return nil
			},
			func() error {
				return errors.New("test error")
			},
		}
		
		err := CurrentLimit(2, workList, false)
		if err == nil {
			t.Error("CurrentLimit() should return error when enableError=false and work fails")
		}
	})
	
	// 测试错误处理 enableError=true
	t.Run("error handling with enableError=true", func(t *testing.T) {
		var results []int
		var mu sync.Mutex
		
		workList := []func() error{
			func() error {
				mu.Lock()
				results = append(results, 1)
				mu.Unlock()
				return nil
			},
			func() error {
				mu.Lock()
				results = append(results, 2)
				mu.Unlock()
				return errors.New("test error")
			},
		}
		
		err := CurrentLimit(2, workList, true)
		if err != nil {
			t.Errorf("CurrentLimit() error = %v, want nil when enableError=true", err)
		}
		
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})
	
	// 测试边界情况：limit <= 0
	t.Run("edge case: limit <= 0", func(t *testing.T) {
		workList := []func() error{
			func() error {
				return nil
			},
		}
		
		err := CurrentLimit(0, workList, false)
		if err != nil {
			t.Errorf("CurrentLimit() error = %v, want nil for limit <= 0", err)
		}
	})
	
	// 测试边界情况：空workList
	t.Run("edge case: empty workList", func(t *testing.T) {
		var workList []func() error
		
		err := CurrentLimit(2, workList, false)
		if err != nil {
			t.Errorf("CurrentLimit() error = %v, want nil for empty workList", err)
		}
	})
}

func TestPaginated(t *testing.T) {
	// 测试正常情况
	t.Run("normal case", func(t *testing.T) {
		var pages []int
		var mu sync.Mutex
		
		err := Paginated(100, 10, func(page int) error {
			mu.Lock()
			pages = append(pages, page)
			mu.Unlock()
			return nil
		})
		
		if err != nil {
			t.Errorf("Paginated() error = %v, want nil", err)
		}
		
		// 100条数据，每页10条，应该有11页 (0-10)
		expectedPages := 11
		if len(pages) != expectedPages {
			t.Errorf("Expected %d pages, got %d", expectedPages, len(pages))
		}
	})
	
	// 测试错误处理
	t.Run("error handling", func(t *testing.T) {
		err := Paginated(100, 10, func(page int) error {
			if page == 5 {
				return errors.New("test error")
			}
			return nil
		})
		
		if err == nil {
			t.Error("Paginated() should return error when work function fails")
		}
	})
	
	// 测试边界情况：total为0
	t.Run("edge case: total is 0", func(t *testing.T) {
		var pages []int
		var mu sync.Mutex
		
		err := Paginated(0, 10, func(page int) error {
			mu.Lock()
			pages = append(pages, page)
			mu.Unlock()
			return nil
		})
		
		if err != nil {
			t.Errorf("Paginated() error = %v, want nil", err)
		}
		
		// total为0时，应该仍然执行一次page 0
		if len(pages) != 1 || pages[0] != 0 {
			t.Errorf("Expected [0] pages, got %v", pages)
		}
	})
}

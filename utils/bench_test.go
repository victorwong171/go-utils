package utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// BenchmarkGetUuid benchmarks UUID generation
func BenchmarkGetUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetUuid()
	}
}

// BenchmarkTernaryOperator benchmarks ternary operator performance
func BenchmarkTernaryOperator(b *testing.B) {
	trueVal := "true"
	falseVal := "false"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := TernaryOperator(i%2 == 0, trueVal, falseVal)
		_ = result
	}
}

// BenchmarkContextCopy benchmarks context copying
func BenchmarkContextCopy(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newCtx := ContextCopy(ctx)
		_ = newCtx
	}
}

// BenchmarkGenerateRandomIntList benchmarks random list generation
func BenchmarkGenerateRandomIntList(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				list := GenerateRandomIntList(size, 1, 100)
				_ = list
			}
		})
	}
}

// BenchmarkMax benchmarks Max function
func BenchmarkMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := Max(i, i+1)
		_ = result
	}
}

// BenchmarkMin benchmarks Min function
func BenchmarkMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := Min(i, i+1)
		_ = result
	}
}

// BenchmarkCurrentLimit benchmarks concurrent processing
func BenchmarkCurrentLimit(b *testing.B) {
	workList := make([]func() error, 100)
	for i := range workList {
		workList[i] = func() error {
			time.Sleep(1 * time.Microsecond)
			return nil
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CurrentLimit(10, workList, false)
	}
}

// BenchmarkPaginated benchmarks pagination processing
func BenchmarkPaginated(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Paginated(1000, 10, func(page int) error {
			return nil
		})
	}
}

// BenchmarkContain benchmarks slice containment check
func BenchmarkContain(b *testing.B) {
	list := make([]int, 1000)
	for i := range list {
		list[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Contain(list, i%1000)
	}
}

// BenchmarkGetItem benchmarks slice item retrieval
func BenchmarkGetItem(b *testing.B) {
	list := make([]int, 1000)
	for i := range list {
		list[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetItem(i%1000, list)
	}
}

// BenchmarkRemoveItemByValue benchmarks item removal
func BenchmarkRemoveItemByValue(b *testing.B) {
	list := make([]int, 1000)
	for i := range list {
		list[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newList := RemoveItemByValue(list, i%1000)
		_ = newList
	}
}

// BenchmarkInsertItems benchmarks item insertion
func BenchmarkInsertItems(b *testing.B) {
	list := make([]int, 1000)
	for i := range list {
		list[i] = i
	}

	items := []int{999, 1000, 1001}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newList, err := InsertItems(list, i%1000, items...)
		if err != nil {
			b.Fatal(err)
		}
		_ = newList
	}
}

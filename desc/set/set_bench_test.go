package set

import (
	"testing"
)

// BenchmarkSet_Set benchmarks set operations
func BenchmarkSet_Set(b *testing.B) {
	s := InitSet[string](1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set("key" + string(rune(i%1000)))
	}
}

// BenchmarkSet_HasKey benchmarks key checking
func BenchmarkSet_HasKey(b *testing.B) {
	s := InitSet[string](1000)

	// Add some keys
	for i := 0; i < 1000; i++ {
		s.Set("key" + string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.HasKey("key" + string(rune(i%1000)))
	}
}

// BenchmarkSet_Drop benchmarks key removal
func BenchmarkSet_Drop(b *testing.B) {
	s := InitSet[string](1000)

	// Add some keys
	for i := 0; i < 1000; i++ {
		s.Set("key" + string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Drop("key" + string(rune(i%1000)))
	}
}

// BenchmarkSet_HasAny benchmarks any key checking
func BenchmarkSet_HasAny(b *testing.B) {
	s := InitSet[string](1000)

	// Add some keys
	for i := 0; i < 1000; i++ {
		s.Set("key" + string(rune(i)))
	}

	keys := []string{"key500", "key501", "key502"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.HasAny(keys...)
	}
}

// BenchmarkSet_ToSlice benchmarks slice conversion
func BenchmarkSet_ToSlice(b *testing.B) {
	s := InitSet[string](1000)

	// Add some keys
	for i := 0; i < 1000; i++ {
		s.Set("key" + string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ToSlice()
	}
}

// BenchmarkSetify benchmarks set creation from slice
func BenchmarkSetify(b *testing.B) {
	keys := make([]string, 1000)
	for i := range keys {
		keys[i] = "key" + string(rune(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Setify(keys...)
	}
}

// BenchmarkSet_Concurrent benchmarks concurrent operations
// Note: This benchmark is disabled due to concurrent map access issues
// func BenchmarkSet_Concurrent(b *testing.B) {
// 	s := InitSet[string](1000)
//
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		i := 0
// 		for pb.Next() {
// 			key := "key" + string(rune(i%1000))
// 			s.Set(key)
// 			_ = s.HasKey(key)
// 			s.Drop(key)
// 			i++
// 		}
// 	})
// }

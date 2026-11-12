package fastnum

import (
	"slices"
	"testing"
)

type number interface {
	int64 | float64
}

var benchSizes = []struct {
	name string
	n    int
}{
	{"8", 8},
	{"64", 64},
	{"1K", 1 << 10},
	{"8K", 8 << 10},
	{"64K", 1 << 16},
	{"16M", 1 << 24},
}

var sink int

// Baseline loops for comparison.

func appendInt64ZerosLoop(dst []int64, n int) []int64 {
	m := len(dst)
	dst = slices.Grow(dst, n)
	clear(dst[m : m+n])
	return dst[:m+n]
}

func appendInt64OnesLoop(dst []int64, n int) []int64 {
	m := len(dst)
	dst = slices.Grow(dst, n)
	t := dst[:m+n]
	for i := m; i < len(t); i++ {
		t[i] = 1
	}
	return t
}

func appendFloat64ZerosLoop(dst []float64, n int) []float64 {
	m := len(dst)
	dst = slices.Grow(dst, n)
	clear(dst[m : m+n])
	return dst[:m+n]
}

func appendFloat64OnesLoop(dst []float64, n int) []float64 {
	m := len(dst)
	dst = slices.Grow(dst, n)
	t := dst[:m+n]
	for i := m; i < len(t); i++ {
		t[i] = 1
	}
	return t
}

func appendInt64ZerosNaive(dst []int64, n int) []int64 {
	for i := 0; i < n; i++ {
		dst = append(dst, 0)
	}
	return dst
}

func appendInt64OnesNaive(dst []int64, n int) []int64 {
	for i := 0; i < n; i++ {
		dst = append(dst, 1)
	}
	return dst
}

func appendFloat64ZerosNaive(dst []float64, n int) []float64 {
	for i := 0; i < n; i++ {
		dst = append(dst, 0)
	}
	return dst
}

func appendFloat64OnesNaive(dst []float64, n int) []float64 {
	for i := 0; i < n; i++ {
		dst = append(dst, 1)
	}
	return dst
}

func isInt64ZerosLoop(a []int64) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

func isInt64OnesLoop(a []int64) bool {
	for _, v := range a {
		if v != 1 {
			return false
		}
	}
	return true
}

func isFloat64ZerosLoop(a []float64) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

func isFloat64OnesLoop(a []float64) bool {
	for _, v := range a {
		if v != 1 {
			return false
		}
	}
	return true
}

// Benchmark harnesses.

func benchmarkAppend[T number](b *testing.B, impl, loop, naive func([]T, int) []T) {
	for _, sz := range benchSizes {
		b.Run(sz.name, func(b *testing.B) {
			buf := make([]T, sz.n)
			b.Run("impl", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r := impl(buf[:0], sz.n)
					sink += len(r)
				}
			})
			b.Run("loop", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r := loop(buf[:0], sz.n)
					sink += len(r)
				}
			})
			b.Run("naive", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r := naive(buf[:0], sz.n)
					sink += len(r)
				}
			})
		})
	}
}

func benchmarkIs[T number](b *testing.B, impl, loop func([]T) bool, val T) {
	for _, sz := range benchSizes {
		b.Run(sz.name, func(b *testing.B) {
			a := make([]T, sz.n)
			for i := range a {
				a[i] = val
			}
			b.Run("impl", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if impl(a) {
						sink++
					}
				}
			})
			b.Run("loop", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if loop(a) {
						sink++
					}
				}
			})
		})
	}
}

func BenchmarkAppendInt64Zeros(b *testing.B) {
	benchmarkAppend(b, AppendInt64Zeros, appendInt64ZerosLoop, appendInt64ZerosNaive)
}

func BenchmarkAppendInt64Ones(b *testing.B) {
	benchmarkAppend(b, AppendInt64Ones, appendInt64OnesLoop, appendInt64OnesNaive)
}

func BenchmarkAppendFloat64Zeros(b *testing.B) {
	benchmarkAppend(b, AppendFloat64Zeros, appendFloat64ZerosLoop, appendFloat64ZerosNaive)
}

func BenchmarkAppendFloat64Ones(b *testing.B) {
	benchmarkAppend(b, AppendFloat64Ones, appendFloat64OnesLoop, appendFloat64OnesNaive)
}

func BenchmarkIsInt64Zeros(b *testing.B) {
	benchmarkIs(b, IsInt64Zeros, isInt64ZerosLoop, int64(0))
}

func BenchmarkIsInt64Ones(b *testing.B) {
	benchmarkIs(b, IsInt64Ones, isInt64OnesLoop, int64(1))
}

func BenchmarkIsFloat64Zeros(b *testing.B) {
	benchmarkIs(b, IsFloat64Zeros, isFloat64ZerosLoop, 0.0)
}

func BenchmarkIsFloat64Ones(b *testing.B) {
	benchmarkIs(b, IsFloat64Ones, isFloat64OnesLoop, 1.0)
}

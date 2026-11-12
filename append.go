package fastnum

import "slices"

const (
	seedElems  = 1 << 10 // 8 KiB; must stay L1-resident
	smallElems = 16
)

// AppendInt64Zeros appends items zeros to dst and returns the result.
//
// It is faster than the corresponding loop.
func AppendInt64Zeros(dst []int64, items int) []int64 {
	return appendInt64Data(dst, items, int64Zeros[:])
}

// AppendInt64Ones appends items ones to dst and returns the result.
//
// It is faster than the corresponding loop.
func AppendInt64Ones(dst []int64, items int) []int64 {
	return appendOnesData(dst, items, int64Ones[:])
}

// AppendFloat64Zeros appends items zeros to dst and returns the result.
//
// It is faster than the corresponding loop.
func AppendFloat64Zeros(dst []float64, items int) []float64 {
	return appendFloat64Data(dst, items, float64Zeros[:])
}

// AppendFloat64Ones appends items ones to dst and returns the result.
//
// It is faster than the corresponding loop.
func AppendFloat64Ones(dst []float64, items int) []float64 {
	if items > 0 && items <= smallElems {
		return append(dst, float64Ones[:items]...)
	}
	return appendOnesData(dst, items, float64Ones[:])
}

func appendInt64Data(dst []int64, items int, src []int64) []int64 {
	for items > 0 {
		n := min(len(src), items)
		dst = append(dst, src[:n]...)
		items -= n
	}
	return dst
}

func appendFloat64Data(dst []float64, items int, src []float64) []float64 {
	for items > 0 {
		n := min(len(src), items)
		dst = append(dst, src[:n]...)
		items -= n
	}
	return dst
}

// appendOnesData appends items ones to dst. Small runs use a plain loop.
// Up to the seed table size it is a single append from the table. Larger runs
// copy an L1-sized seed prefix once and then keep copying from that hot prefix,
// so the copy source never leaves L1.
func appendOnesData[T int64 | float64](dst []T, items int, seed []T) []T {
	if items <= 0 {
		return dst
	}
	if items <= smallElems {
		m := len(dst)
		dst = slices.Grow(dst, items)
		t := dst[:m+items]
		for i := m; i < len(t); i++ {
			t[i] = 1
		}
		return t
	}
	if items <= len(seed) {
		return append(dst, seed[:items]...)
	}
	m := len(dst)
	dst = slices.Grow(dst, items)
	t := dst[:m+items]
	s := t[m:]
	n := copy(s, seed)
	for f := n; f < len(s); f += n {
		copy(s[f:], s[:n])
	}
	return t
}

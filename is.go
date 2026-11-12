package fastnum

import (
	"bytes"
	"math"
	"unsafe"
)

// IsInt64Zeros checks whether a contains only zeros.
func IsInt64Zeros(a []int64) bool {
	return isInt64SameValue(a, 0)
}

// IsInt64Ones checks whether a contains only ones.
func IsInt64Ones(a []int64) bool {
	return isInt64SameValue(a, 1)
}

// IsFloat64Zeros checks whether a contains only zeros.
func IsFloat64Zeros(a []float64) bool {
	if len(a) == 0 {
		return true
	}
	if math.Float64bits(a[0]) != 0 {
		return false
	}
	return isFloat64SameBits(a)
}

// IsFloat64Ones checks whether a contains only ones.
func IsFloat64Ones(a []float64) bool {
	if len(a) == 0 {
		return true
	}
	if math.Float64bits(a[0]) != math.Float64bits(1) {
		return false
	}
	return isFloat64SameBits(a)
}

// isInt64SameValue reports whether a consists of the same v values.
//
// Comparing the slice bytes against themselves shifted by one element
// proves that all elements are equal; a[0] then pins the exact value.
// Unlike a reference table, this reads only one memory stream.
func isInt64SameValue(a []int64, v int64) bool {
	if len(a) == 0 {
		return true
	}
	if a[0] != v {
		return false
	}
	if len(a) <= smallElems {
		for _, x := range a[1:] {
			if x != v {
				return false
			}
		}
		return true
	}
	b := int64ToByteSlice(a)
	return bytes.Equal(b[8:], b[:len(b)-8])
}

// isFloat64SameBits reports whether all elements of a share the same
// bit pattern as a[0].
func isFloat64SameBits(a []float64) bool {
	if len(a) <= smallElems {
		bits := math.Float64bits(a[0])
		for _, x := range a[1:] {
			if math.Float64bits(x) != bits {
				return false
			}
		}
		return true
	}
	b := float64ToByteSlice(a)
	return bytes.Equal(b[8:], b[:len(b)-8])
}

func int64ToByteSlice(a []int64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(a))), len(a)*8)
}

func float64ToByteSlice(a []float64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(a))), len(a)*8)
}

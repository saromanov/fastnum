package fastnum

import (
	"runtime"
	"sync"
	"testing"
)

func fillChunk[T int64 | float64](s []T, seed []T) {
	n := copy(s, seed)
	for f := n; f < len(s); f += n {
		copy(s[f:], s[:n])
	}
}

// BenchmarkParallelFillOnes compares serial and parallel ones-fill at
// sizes where the fill is memory-bandwidth bound.
func BenchmarkParallelFillOnes(b *testing.B) {
	const n = 1 << 24
	seed := int64Ones[:]
	b.Run("serial", func(b *testing.B) {
		buf := make([]int64, n)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fillChunk(buf, seed)
		}
	})
	for _, nw := range []int{2, 4, runtime.NumCPU()} {
		b.Run("par"+itoa(nw), func(b *testing.B) {
			buf := make([]int64, n)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				ch := (n + nw - 1) / nw
				ch -= ch % len(seed)
				wg.Add(nw)
				for w := 0; w < nw; w++ {
					go func(w int) {
						defer wg.Done()
						lo, hi := w*ch, min((w+1)*ch, n)
						if lo < hi {
							fillChunk(buf[lo:hi], seed)
						}
					}(w)
				}
				wg.Wait()
			}
		})
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	p := len(b)
	for i > 0 {
		p--
		b[p] = byte('0' + i%10)
		i /= 10
	}
	return string(b[p:])
}

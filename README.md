# fastnum

[![CI](https://github.com/saromanov/fastnum/actions/workflows/ci.yml/badge.svg)](https://github.com/saromanov/fastnum/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/saromanov/fastnum.svg)](https://pkg.go.dev/github.com/saromanov/fastnum)
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/fastnum)](https://goreportcard.com/report/github.com/saromanov/fastnum)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Fast fill and constant-value checks for numeric slices in Go.

`fastnum` provides optimized helpers for everyday operations on `[]int64`
and `[]float64`: appending zeros or ones, and testing whether a slice
consists entirely of zeros or ones. The functions are designed to stay
within a single CPU core's caches and to avoid allocations when the
destination buffer has enough capacity.

The API mirrors [`lib/fastnum`](https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/lib/fastnum)
of [VictoriaMetrics](https://victoriametrics.com/) (Apache-2.0), from which
this package is derived, so it can be used as a drop-in replacement.

## Features

- `AppendInt64Zeros`, `AppendInt64Ones`, `AppendFloat64Zeros`, `AppendFloat64Ones` — append-style fills following the standard library convention, so buffers can be reused across calls.
- `IsInt64Zeros`, `IsInt64Ones`, `IsFloat64Zeros`, `IsFloat64Ones` — constant-slice checks with early exit on the first mismatch.
- Zero allocations when the destination has sufficient capacity.
- Correct handling of `-0.0` and `NaN` bit patterns in float64 checks.
- Standard library only, no external dependencies.

## Installation

```sh
go get github.com/saromanov/fastnum
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/saromanov/fastnum"
)

func main() {
	// Reuse dst across calls to avoid allocations.
	var dst []int64
	dst = fastnum.AppendInt64Ones(dst[:0], 4)
	fmt.Println(dst) // [1 1 1 1]

	fmt.Println(fastnum.IsInt64Zeros([]int64{0, 0, 0})) // true
	fmt.Println(fastnum.IsInt64Zeros([]int64{0, 1, 0})) // false
}
```

## How it works

- **L1-resident seed table for ones.** The ones seed table is 8 KiB (down
  from 64 KiB in the original implementation), so the copy source never
  leaves L1.
- **Prefix-doubling fill.** Slices larger than the seed table are filled by
  copying the seed once and then repeatedly from the freshly written hot
  prefix, turning a DRAM-streaming operation into an L1-bound one.
- **Shifted-byte comparison.** A constant check compares the slice bytes
  against themselves shifted by one element (`bytes.Equal(b[8:], b[:len(b)-8])`):
  a single read stream proves all elements are equal, and the first element
  pins the exact value. float64 checks compare bit patterns, preserving
  `-0`/`NaN` semantics.
- **Early exit.** The first element is validated before any full scan, so
  mismatches return in O(1).

## Benchmarks

Environment: Intel Core i7-9750H (6C/12T), Go 1.27.3, linux/amd64,
`GOAMD64=v1`, destination buffer preallocated (0 allocs in all rows).
Medians of 3 runs; small-size numbers vary ±10–20% between runs.
`fastnum` vs a plain loop (`slices.Grow` + fill/`clear`):

| Operation        | 1K             | 8K             | 64K              | 16M              |
|------------------|---------------:|---------------:|-----------------:|-----------------:|
| AppendInt64Ones  | 79 / 540 ns    | 1.0 / 4.4 µs   | 13.3 / 34.8 µs   | 11.5 / 11.6 ms   |
| AppendFloat64Ones| 95 / 570 ns    | 1.0 / 4.5 µs   | 13.6 / 34.3 µs   | 11.3 / 10.7 ms   |
| IsInt64Zeros     | 125 / 305 ns   | 1.2 / 2.3 µs   | 14.5 / 22.5 µs   | 8.3 / 10.0 ms    |
| IsFloat64Zeros   | 125 / 600 ns   | 1.3 / 4.9 µs   | 11.8 / 40.4 µs   | 7.5 / 12.9 ms    |

Ones fills and constant checks win from ~1K elements up to L3-sized slices;
at DRAM-sized slices everything is memory-bound and on par. Zeros appends
are on par with a `clear`-based loop (the runtime's `clear` uses wider
stores, which wins at 16M) and significantly faster than a naive
`append` in a loop.

Full methodology, raw results and the research notes (including why
non-temporal stores are the remaining upside for huge ones fills) live in
[docs/BENCHMARKS.md](docs/BENCHMARKS.md) and [docs/RESEARCH.md](docs/RESEARCH.md).
Reproduce with:

```sh
go test -run '^$' -bench . -benchmem -count=6 -benchtime=500ms
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and pull requests are
welcome.

## License

[MIT](LICENSE). Portions derived from VictoriaMetrics `lib/fastnum`,
© VictoriaMetrics, Apache-2.0.

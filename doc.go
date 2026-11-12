// Package fastnum provides fast functions for filling int64 and float64
// slices with zeros or ones and for checking whether such slices consist
// entirely of one repeated value.
//
// All append-style functions follow the standard library convention: they
// append to the provided dst and return the resulting slice, so buffers can
// be reused across calls without extra allocations.
//
// The package is derived from the lib/fastnum package of VictoriaMetrics
// (Apache-2.0) with additional optimizations: L1-resident seed tables for
// ones, prefix-doubling fill for large slices, and single-stream
// shifted-byte comparison for constant checks.
package fastnum

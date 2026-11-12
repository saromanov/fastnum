# Changelog

## v0.1.0 (2026-11-12)

- Initial release.
- `AppendInt64Zeros`, `AppendInt64Ones`, `AppendFloat64Zeros`,
  `AppendFloat64Ones`: L1-resident 8 KiB ones seed table, single-append fast
  path up to the seed size, prefix-doubling fill for larger runs.
- `IsInt64Zeros`, `IsInt64Ones`, `IsFloat64Zeros`, `IsFloat64Ones`:
  single-stream shifted-byte comparison with early exit on the first
  mismatch; float64 checks compare bit patterns, preserving `-0`/`NaN`
  semantics.

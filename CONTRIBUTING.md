# Contributing

Contributions are welcome. Please open an issue before starting work on
larger changes.

## Prerequisites

- Go 1.26 or newer.

## Local checks

```sh
gofmt -l .
go vet ./...
go test -race -count=1 ./...
```

CI additionally runs `golangci-lint` and `govulncheck`.

## Benchmarks

Performance-sensitive changes should come with benchmark results following
the methodology in [docs/BENCHMARKS.md](docs/BENCHMARKS.md):

```sh
go test -run '^$' -bench . -benchmem -count=6 -benchtime=500ms | tee out.txt
benchstat out.txt
```

Report the machine, Go version, `GOAMD64` level, and medians.

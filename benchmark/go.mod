module github.com/nxtcoder17/fastlog/benchmark

go 1.24.1

require github.com/rs/zerolog v1.34.0

require (
	github.com/nxtcoder17/fastlog v0.0.0
	github.com/phuslu/log v1.0.115
	go.uber.org/zap v1.27.0
)

require (
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/nxtcoder17/fastlog v0.0.0 => ../.

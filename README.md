# fastlog - A minimal, faster logger for Go

fastlog provides a high-performance logging library with a clean fluent API. It supports multiple output formats and is significantly faster than popular Go loggers.

## Features

- **Fast** - 2-3x faster than zerolog, zap, and phuslu-log
- **Clean API** - Fluent builder pattern
- **Multiple formats** - Console, JSON, Logfmt
- **Thread-safe** - Built-in sync writer
- **Memory efficient** - Buffer pool, minimal allocations
- **slog compatible** - Returns `*slog.Logger` via `.Slog()`

## Getting Started

```go
import "github.com/nxtcoder17/fastlog"

func main() {
    // Quick start with default settings
    logger := fastlog.Console()
    logger.Info("hello world")

    // Fluent API for customization
    logger := fastlog.New().
        Writer(os.Stderr).
        Colors(true).
        Timestamp(true).
        Caller(true).
        Logfmt()
    
    logger.Debug("debug message", "key", "value")
    logger.Info("info message", "user", "john", "age", 30)
    logger.Warn("warning message", "err", fmt.Errorf("something failed"))
    logger.Error("error message", "err", fmt.Errorf("failed to connect"))

    // slog compatibility
    slog := logger.Slog()
    slog.Info("using standard slog interface")

    // Clone logger with new settings
    cloned := logger.Clone().DebugMode(true).JSON()
}
```

## API Options

| Method                  | Description                            |
| --------                | -------------                          |
| `Writer(io.Writer)`     | Set output writer (default: os.Stderr) |
| `Console()`             | Console format output                  |
| `JSON()`                | JSON format output                     |
| `Logfmt()`              | Logfmt format output                   |
| `Colors(bool)`          | Enable/disable colors (default: true)  |
| `Timestamp(bool)`       | Show/hide timestamp (default: true)    |
| `Caller(bool)`          | Show/hide caller (default: true)       |
| `LogLevel(slog.Level)`  | Set minimum log level                  |
| `DebugMode(bool)`       | Shortcut to enable debug logging       |
| `Prefix(string)`        | Set log prefix                         |
| `SkipCallerFrames(int)` | Skip additional caller frames          |
| --------                | -------------                          |

## Quick Start Functions

For simple use cases, use the package-level functions:

```go
fastlog.Console()  // Console format with defaults
fastlog.JSON()     // JSON format with defaults
fastlog.Logfmt()   // Logfmt format with defaults
```

## Example Output

### Console Format
```
2024-01-15T10:30:45Z | /path/to/file.go:42 | INFO | hello world | user=john age=30
```

### JSON Format
```
{"timestamp":"2024-01-15T10:30:45Z","level":"INFO","message":"hello world","user":"john","age":30}
```

### Logfmt Format
```
timestamp=2024-01-15T10:30:45Z level=INFO message=hello world user=john age=30
```

## Benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/nxtcoder17/fastlog/benchmark
cpu: AMD Ryzen 9 6900HS with Radeon Graphics
BenchmarkFastlog_console_withoutCaller-16         	  713667	      1688 ns/op	     272 B/op	      16 allocs/op
BenchmarkFastlog_console_withCaller-16            	  481687	      2348 ns/op	     520 B/op	      18 allocs/op
BenchmarkFastlog_console_slog_withoutCaller-16    	  437119	      2819 ns/op	     600 B/op	      28 allocs/op
BenchmarkFastlog_console_slog_withCaller-16       	  291745	      4001 ns/op	     849 B/op	      30 allocs/op
BenchmarkFastlog_logfmt_withoutCaller-16          	  767406	      1596 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_logfmt_withCaller-16             	  500019	      2224 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_logfmt_slog_withoutCaller-16     	  421028	      2768 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_slog_withCaller-16        	  302834	      3956 ns/op	     849 B/op	      31 allocs/op
BenchmarkFastlog_json_withoutCaller-16            	  626450	      1815 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_json_withCaller-16               	  499382	      2398 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_json_slog_withoutCaller-16       	  446340	      2621 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_json_slog_withCaller-16          	  295039	      3726 ns/op	     849 B/op	      31 allocs/op
BenchmarkPhusluLog_withoutCaller-16               	  363284	      3259 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_withCaller-16                  	  342709	      3437 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_slog_withoutCaller-16          	  304657	      3838 ns/op	    1106 B/op	      29 allocs/op
BenchmarkPhusluLog_slog_withCaller-16             	  283930	      4176 ns/op	    1106 B/op	      29 allocs/op
BenchmarkSlog_JSON-16                             	  247533	      4817 ns/op	    1738 B/op	      40 allocs/op
BenchmarkSlog_Info_With_Caller-16                 	  207474	      5579 ns/op	    2323 B/op	      46 allocs/op
BenchmarkSlog_Text-16                             	  158947	      6557 ns/op	    1442 B/op	      34 allocs/op
BenchmarkSlog_Text_With_Caller-16                 	  169699	      7488 ns/op	    1834 B/op	      39 allocs/op
BenchmarkZap_SugarInfo-16                         	  235269	      4915 ns/op	    2593 B/op	      32 allocs/op
BenchmarkZerolog_Info-16                          	  291512	      3953 ns/op	    2018 B/op	      46 allocs/op
```

## Performance Comparison (logfmt without caller)

| Logger     | ns/op   | Relative Speed    |
| --------   | ------- | ----------------- |
| fastlog    | 1,570   | 1x (baseline)     |
| phuslu-log | 3,178   | 2.0x slower       |
| zerolog    | 3,910   | 2.5x slower       |
| zap        | 4,831   | 3.1x slower       |
| slog       | 4,747   | 3.0x slower       |

## Why is fastlog fast?

- **Buffer pooling** - Reuses byte buffers via `sync.Pool`
- **Direct byte manipulation** - No `fmt.Sprintf` or `reflect`
- **Type-specific serialization** - Fast path for common types
- **Zero-copy string escaping** - Optimized JSON string handling
- **Minimal allocations** - Pre-allocated buffers, inline operations

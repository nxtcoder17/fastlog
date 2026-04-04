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
BenchmarkFastlog_console_withoutCaller-16         	  750530	      1687 ns/op	     272 B/op	      16 allocs/op
BenchmarkFastlog_console_withCaller-16            	  510012	      2367 ns/op	     520 B/op	      18 allocs/op
BenchmarkFastlog_console_slog_withoutCaller-16    	  395643	      2850 ns/op	     600 B/op	      28 allocs/op
BenchmarkFastlog_console_slog_withCaller-16       	  292074	      3850 ns/op	     849 B/op	      30 allocs/op
BenchmarkFastlog_logfmt_withoutCaller-16          	  754854	      1554 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_logfmt_withCaller-16             	  507148	      2136 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_logfmt_slog_withoutCaller-16     	  442084	      2703 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_slog_withCaller-16        	  294006	      3856 ns/op	     849 B/op	      31 allocs/op
BenchmarkFastlog_json_withoutCaller-16            	  631666	      1785 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_json_withCaller-16               	  511449	      2361 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_json_slog_withoutCaller-16       	  492972	      2537 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_json_slog_withCaller-16          	  339217	      3643 ns/op	     849 B/op	      31 allocs/op
BenchmarkPhusluLog_withoutCaller-16               	  355879	      3118 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_withCaller-16                  	  344761	      3493 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_slog_withoutCaller-16          	  309410	      3848 ns/op	    1106 B/op	      29 allocs/op
BenchmarkPhusluLog_slog_withCaller-16             	  295239	      3976 ns/op	    1106 B/op	      29 allocs/op
BenchmarkSlog_JSON-16                             	  249128	      4721 ns/op	    1738 B/op	      40 allocs/op
BenchmarkSlog_Info_With_Caller-16                 	  194593	      5881 ns/op	    2323 B/op	      46 allocs/op
BenchmarkSlog_Text-16                             	  165500	      6754 ns/op	    1442 B/op	      34 allocs/op
BenchmarkSlog_Text_With_Caller-16                 	  158253	      7313 ns/op	    1850 B/op	      39 allocs/op
BenchmarkZap_SugarInfo-16                         	  234570	      4962 ns/op	    2593 B/op	      32 allocs/op
BenchmarkZerolog_Info-16                          	  287758	      3950 ns/op	    2018 B/op	      46 allocs/op
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

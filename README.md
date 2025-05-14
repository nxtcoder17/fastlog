## fastlog is a minimal, faster logger for go

fastlog provides a slog Handler implementation in logfmt and json formats.

In my initial benchmarks, it tends to perform better than some other popular loggers

## Getting Started

```go
import (
    "github.com/nxtcoder17/fastlog"
)

func main() {
    // you can also choose other formats `Logfmt` and `JSON`
    logfmt := fastlog.New(fastlog.Options{Format: fastlog.ConsoleFormat, ShowCaller: true, EnableColors: true})
}
```

## Output

![Example](https://github.com/user-attachments/assets/a89ab883-fea1-462c-9c67-dde889c4b1e0)

## benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/nxtcoder17/fastlog/benchmark
cpu: AMD Ryzen 9 6900HS with Radeon Graphics        
BenchmarkFastlog_console_withoutCaller-16         	  791881	      1458 ns/op	     224 B/op	      14 allocs/op
BenchmarkFastlog_console_withCaller-16            	  589026	      2106 ns/op	     472 B/op	      16 allocs/op
BenchmarkFastlog_console_slog_withoutCaller-16    	  464258	      2486 ns/op	     552 B/op	      26 allocs/op
BenchmarkFastlog_console_slog_withCaller-16       	  348403	      3400 ns/op	     801 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_withoutCaller-16          	  779101	      1501 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_logfmt_withCaller-16             	  549084	      2122 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_logfmt_slog_withoutCaller-16     	  425734	      2706 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_slog_withCaller-16        	  319371	      3716 ns/op	     849 B/op	      31 allocs/op
BenchmarkFastlog_json_withoutCaller-16            	  637552	      1836 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_json_withCaller-16               	  501534	      2411 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_json_slog_withoutCaller-16       	  452623	      2583 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_json_slog_withCaller-16          	  343165	      3472 ns/op	     849 B/op	      31 allocs/op
BenchmarkPhusluLog_withoutCaller-16               	  353538	      3168 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_withCaller-16                  	  336570	      3496 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_slog_withoutCaller-16          	  295429	      3915 ns/op	    1106 B/op	      29 allocs/op
BenchmarkPhusluLog_slog_withCaller-16             	  268814	      4015 ns/op	    1106 B/op	      29 allocs/op
BenchmarkSlog_Info-16                             	  235711	      4765 ns/op	    1738 B/op	      40 allocs/op
BenchmarkZap_SugarInfo-16                         	  222519	      5299 ns/op	    2594 B/op	      32 allocs/op
BenchmarkZerolog_Info-16                          	  262026	      4326 ns/op	    2643 B/op	      48 allocs/op
PASS
ok  	github.com/nxtcoder17/fastlog/benchmark	23.587s
```

![Benchmark Results](https://github.com/user-attachments/assets/78d3dadf-4009-42a4-899e-332533f7b109)

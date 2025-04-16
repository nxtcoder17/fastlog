## fastlog is a minimal, faster logger for go

fastlog provides a slog Handler implementation in logfmt and json formats.

In my initial benchmarks, it tends to perform better than some other popular loggers

## benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/nxtcoder17/fastlog/benchmark
cpu: AMD Ryzen 9 6900HS with Radeon Graphics        
BenchmarkFastlog_console_withoutCaller-16         	  698872	      1532 ns/op	     224 B/op	      14 allocs/op
BenchmarkFastlog_console_withCaller-16            	  567091	      2053 ns/op	     472 B/op	      16 allocs/op
BenchmarkFastlog_console_slog_withoutCaller-16    	  475928	      2290 ns/op	     552 B/op	      26 allocs/op
BenchmarkFastlog_console_slog_withCaller-16       	  317322	      3199 ns/op	     801 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_withoutCaller-16          	  717027	      1601 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_logfmt_withCaller-16             	  534128	      2145 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_logfmt_slog_withoutCaller-16     	  447391	      2490 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_logfmt_slog_withCaller-16        	  315032	      3518 ns/op	     849 B/op	      31 allocs/op
BenchmarkFastlog_json_withoutCaller-16            	  639082	      1723 ns/op	     256 B/op	      16 allocs/op
BenchmarkFastlog_json_withCaller-16               	  524022	      2265 ns/op	     520 B/op	      19 allocs/op
BenchmarkFastlog_json_slog_withoutCaller-16       	  447430	      2577 ns/op	     584 B/op	      28 allocs/op
BenchmarkFastlog_json_slog_withCaller-16          	  338354	      3501 ns/op	     849 B/op	      31 allocs/op
BenchmarkPhusluLog_withoutCaller-16               	  366986	      3125 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_withCaller-16                  	  357410	      3374 ns/op	     945 B/op	      28 allocs/op
BenchmarkPhusluLog_slog_withoutCaller-16          	  288595	      3776 ns/op	    1106 B/op	      29 allocs/op
BenchmarkPhusluLog_slog_withCaller-16             	  293012	      3902 ns/op	    1106 B/op	      29 allocs/op
BenchmarkSlog_Info-16                             	  245907	      4685 ns/op	    1738 B/op	      40 allocs/op
BenchmarkZap_SugarInfo-16                         	  220171	      5246 ns/op	    2593 B/op	      32 allocs/op
BenchmarkZerolog_Info-16                          	  299310	      3924 ns/op	    2018 B/op	      46 allocs/op
PASS
ok  	github.com/nxtcoder17/fastlog/benchmark	23.118s
```

![Benchmark Results](https://github.com/user-attachments/assets/78d3dadf-4009-42a4-899e-332533f7b109)

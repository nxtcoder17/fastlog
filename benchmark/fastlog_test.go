package benchmark

import (
	"io"
	"testing"

	"github.com/nxtcoder17/fastlog"
)

func BenchmarkFastlog_console_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.ConsoleFormat, ShowCaller: false})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.ConsoleFormat, ShowCaller: true})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.ConsoleFormat, ShowCaller: false}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.ConsoleFormat, ShowCaller: true}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.LogfmtFormat, ShowCaller: false})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.LogfmtFormat, ShowCaller: true})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.LogfmtFormat, ShowCaller: false}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.LogfmtFormat, ShowCaller: true}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.JSONFormat, ShowCaller: false})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.JSONFormat, ShowCaller: true})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.JSONFormat, ShowCaller: false}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.Options{Writer: io.Discard, Format: fastlog.JSONFormat, ShowCaller: true}).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

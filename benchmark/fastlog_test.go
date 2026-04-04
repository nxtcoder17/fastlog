package benchmark

import (
	"io"
	"testing"

	"github.com/nxtcoder17/fastlog"
)

func BenchmarkFastlog_console_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).Console()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Console()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).Console().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Console().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).Logfmt()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Logfmt()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).Logfmt().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Logfmt().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).JSON()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).JSON()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).Caller(false).JSON().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withCaller(b *testing.B) {
	logger := fastlog.New().Writer(io.Discard).JSON().Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

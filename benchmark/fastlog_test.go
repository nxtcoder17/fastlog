package benchmark

import (
	"io"
	"testing"

	"github.com/nxtcoder17/fastlog"
)

func BenchmarkFastlog_console_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Console(), fastlog.WithoutCaller())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Console())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Console(), fastlog.WithoutCaller()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_console_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Console()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Logfmt(), fastlog.WithoutCaller())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Logfmt())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Logfmt(), fastlog.WithoutCaller()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_logfmt_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Logfmt()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Json(), fastlog.WithoutCaller())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Json())
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withoutCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Json(), fastlog.WithoutCaller()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkFastlog_json_slog_withCaller(b *testing.B) {
	logger := fastlog.New(fastlog.WithWriter(io.Discard), fastlog.Json()).Slog()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

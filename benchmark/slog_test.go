package benchmark

import (
	"io"
	"log/slog"
	"testing"
)

func BenchmarkSlog_JSON(b *testing.B) {
	handler := slog.NewJSONHandler(io.Discard, nil)
	logger := slog.New(handler)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkSlog_Info_With_Caller(b *testing.B) {
	handler := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(handler)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkSlog_Text(b *testing.B) {
	handler := slog.NewTextHandler(io.Discard, nil)
	logger := slog.New(handler)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

func BenchmarkSlog_Text_With_Caller(b *testing.B) {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(handler)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

package benchmark

import (
	"io"
	"testing"

	"github.com/phuslu/log"
)

func BenchmarkPhusluLog_withoutCaller(b *testing.B) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: io.Discard},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info().KeysAndValues(attrs...).Msg("hello world")
	}
}

func BenchmarkPhusluLog_withCaller(b *testing.B) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: io.Discard},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info().KeysAndValues(attrs...).Msg("hello world")
	}
}

func BenchmarkPhusluLog_slog_withoutCaller(b *testing.B) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: io.Discard},
	}

	b.ReportAllocs()
	b.ResetTimer()

	slog := logger.Slog()

	for i := 0; i < b.N; i++ {
		slog.Info("hello world", attrs...)
	}
}

func BenchmarkPhusluLog_slog_withCaller(b *testing.B) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: io.Discard},
	}

	b.ReportAllocs()
	b.ResetTimer()

	slog := logger.Slog()

	for i := 0; i < b.N; i++ {
		slog.Info("hello world", attrs...)
	}
}

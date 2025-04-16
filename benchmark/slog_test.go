package benchmark

import (
	"io"
	"log/slog"
	"testing"
)

func BenchmarkSlog_Info(b *testing.B) {
	handler := slog.NewJSONHandler(io.Discard, nil)
	logger := slog.New(handler)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello world", attrs...)
	}
}

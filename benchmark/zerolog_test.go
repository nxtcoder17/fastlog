package benchmark

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func BenchmarkZerolog_Info(b *testing.B) {
	logger := zerolog.New(io.Discard).With().Logger()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ev := logger.Info()
		for j := 1; j < len(attrs); j += 2 {
			ev = ev.Any(attrs[j-1].(string), attrs[j])
		}
		ev.Msg("hello world")
		// logger.Info().
		// 	Str("user", "anshuman").
		// 	Int("id", 42).
		// 	Msg("hello world")
	}
}

package benchmark

import (
	"io"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkZap_SugarInfo(b *testing.B) {
	cfg := zap.NewProductionConfig()
	logger, err := cfg.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			cfg.Level,
		)
	}))
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Sugar().Infow("hello world", attrs...)
	}
}

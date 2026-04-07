package fastlog_test

import (
	"fmt"
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"strings"
	"testing"
)

func TestLogfmt_ValidOutput(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).Logfmt()

	logger.Info("hello", "key", "value", "num", 42)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if !strings.Contains(line, "message=hello") {
		t.Errorf("expected message=hello in: %s", line)
	}
	if !strings.Contains(line, "key=value") {
		t.Errorf("expected key=value in: %s", line)
	}
	if !strings.Contains(line, "num=42") {
		t.Errorf("expected num=42 in: %s", line)
	}
}

func TestLogfmt_AllLevels(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		fn    func(fastlog.Logger, string, ...any)
	}{
		{"Debug", slog.LevelDebug, fastlog.Logger.Debug},
		{"Info", slog.LevelInfo, fastlog.Logger.Info},
		{"Warn", slog.LevelWarn, fastlog.Logger.Warn},
		{"Error", slog.LevelError, fastlog.Logger.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := newCapture()
			logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).LogLevel(tt.level).Logfmt()

			tt.fn(logger, "msg", "k", "v")

			lines := cw.Lines()
			if len(lines) != 1 {
				t.Fatalf("expected 1 line, got %d", len(lines))
			}
			if !strings.Contains(lines[0], fmt.Sprintf("level=%s", tt.level.String())) {
				t.Errorf("expected level=%s in: %s", tt.level.String(), lines[0])
			}
		})
	}
}

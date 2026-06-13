package fastlog_test

import (
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"strings"
	"testing"
)

func TestConsole_OutputContainsMessage(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).Console()

	logger.Info("hello world", "key", "value")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "INFO") {
		t.Errorf("expected INFO in: %s", lines[0])
	}
	if !strings.Contains(lines[0], "hello world") {
		t.Errorf("expected 'hello world' in: %s", lines[0])
	}
	if !strings.Contains(lines[0], "key=value") {
		t.Errorf("expected key=value in: %s", lines[0])
	}
}

func TestConsole_AllLevels(t *testing.T) {
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
			logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).LogLevel(tt.level).Console()

			tt.fn(logger, "msg")

			lines := cw.Lines()
			if len(lines) != 1 {
				t.Fatalf("expected 1 line, got %d", len(lines))
			}
			if !strings.Contains(lines[0], tt.level.String()) {
				t.Errorf("expected %s in: %s", tt.level.String(), lines[0])
			}
		})
	}
}

package fastlog_test

import (
	"encoding/json"
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"strings"
	"testing"
)

func TestJSON_ValidOutput(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("hello", "key", "value", "num", 42)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, lines[0])
	}

	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
	if m["num"] != float64(42) {
		t.Errorf("expected num=42, got %v", m["num"])
	}
}

func TestJSON_AllLevels(t *testing.T) {
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
			logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).LogLevel(tt.level).JSON()

			tt.fn(logger, "msg", "k", "v")

			lines := cw.Lines()
			if len(lines) != 1 {
				t.Fatalf("expected 1 line, got %d", len(lines))
			}

			var m map[string]any
			if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}
			if m["level"] != tt.level.String() {
				t.Errorf("expected level=%s, got %v", tt.level.String(), m["level"])
			}
		})
	}
}

func TestJSON_WithCaller(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(true).Timestamp(false).Colors(false).JSON()

	logger.Info("msg")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	caller, ok := m["caller"]
	if !ok {
		t.Fatal("expected caller field in JSON output")
	}
	if !strings.Contains(caller.(string), "json_test.go") {
		t.Errorf("expected caller to contain test filename, got %v", caller)
	}
}

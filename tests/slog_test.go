package fastlog_test

import (
	"encoding/json"
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"strings"
	"testing"
)

func TestSlog_JSON(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON().Slog()

	logger.Info("hello", "key", "value")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestSlog_WithAttrs(t *testing.T) {
	cw := newCapture()
	base := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON().Slog()

	derived := base.With("base_attr", "base_value")
	derived.Info("msg", "call_attr", "call_value")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	json.Unmarshal([]byte(lines[0]), &m)

	if m["base_attr"] != "base_value" {
		t.Errorf("expected base_attr=base_value, got %v", m["base_attr"])
	}
	if m["call_attr"] != "call_value" {
		t.Errorf("expected call_attr=call_value, got %v", m["call_attr"])
	}
}

func TestSlog_WithGroup(t *testing.T) {
	cw := newCapture()
	base := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON().Slog()

	grouped := base.WithGroup("group").With("key", "value")
	grouped.Info("msg")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "group.key") {
		t.Errorf("expected group.key in: %s", lines[0])
	}
}

func TestSlog_LevelFiltering(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).LogLevel(slog.LevelWarn).JSON().Slog()

	logger.Info("should be suppressed")
	logger.Warn("should appear")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Errorf("expected 1 line (Info suppressed), got %d", len(lines))
	}
	if len(lines) > 0 && !strings.Contains(lines[0], "should appear") {
		t.Errorf("expected 'should appear' in output, got: %s", lines[0])
	}
}

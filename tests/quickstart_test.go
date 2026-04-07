package fastlog_test

import (
	"encoding/json"
	"github.com/nxtcoder17/fastlog"
	"strings"
	"testing"
)

func TestQuickStart_Console(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).Console()

	logger.Info("test")

	if len(cw.Lines()) == 0 {
		t.Error("Console() should produce output")
	}
}

func TestQuickStart_JSON(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("test")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestQuickStart_Logfmt(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).Logfmt()

	logger.Info("test")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "message=test") {
		t.Errorf("expected message=test in: %s", lines[0])
	}
}

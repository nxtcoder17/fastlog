package fastlog_test

import (
	"encoding/json"
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"testing"
)

func TestWith_Immutability(t *testing.T) {
	cwBase := newCapture()
	cwDerived := newCapture()

	base := fastlog.New().Writer(cwBase).Caller(false).Timestamp(false).Colors(false).JSON()

	derived := base.With("extra", "attr")

	derivedIsolated := derived.Clone().Writer(cwDerived).JSON()

	base.Info("base msg")
	derivedIsolated.Info("derived msg")

	linesBase := cwBase.Lines()
	if len(linesBase) != 1 {
		t.Fatalf("expected 1 line from base, got %d", len(linesBase))
	}
	var mBase map[string]any
	json.Unmarshal([]byte(linesBase[0]), &mBase)
	if _, ok := mBase["extra"]; ok {
		t.Error("base logger should not have 'extra' attr after With()")
	}

	linesDerived := cwDerived.Lines()
	if len(linesDerived) != 1 {
		t.Fatalf("expected 1 line from derived, got %d", len(linesDerived))
	}
	var mDerived map[string]any
	json.Unmarshal([]byte(linesDerived[0]), &mDerived)
	if mDerived["extra"] != "attr" {
		t.Errorf("derived logger should have extra=attr, got %v", mDerived["extra"])
	}
}

func TestWith_Chaining(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON().
		With("a", "1").
		With("b", "2")

	logger.Info("msg")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	json.Unmarshal([]byte(lines[0]), &m)

	if m["a"] != "1" {
		t.Errorf("expected a=1, got %v", m["a"])
	}
	if m["b"] != "2" {
		t.Errorf("expected b=2, got %v", m["b"])
	}
}

func TestClone_Independent(t *testing.T) {
	cw1 := newCapture()
	cw2 := newCapture()

	logger := fastlog.New().Writer(cw1).Caller(false).Timestamp(false).Colors(false).JSON()
	clonedBuilder := logger.Clone().Writer(cw2)
	cloned := clonedBuilder.JSON()

	logger.Info("original")
	cloned.Info("cloned")

	if len(cw1.Lines()) != 1 {
		t.Errorf("expected 1 line from original, got %d", len(cw1.Lines()))
	}
	if len(cw2.Lines()) != 1 {
		t.Errorf("expected 1 line from cloned, got %d", len(cw2.Lines()))
	}
}

func TestClone_PreservesSettings(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(true).Colors(false).LogLevel(slog.LevelDebug).JSON()
	cloned := logger.Clone()

	clonedWriter := newCapture()
	clonedLogger := cloned.Writer(clonedWriter).JSON()

	clonedLogger.Debug("debug msg")

	lines := clonedWriter.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	json.Unmarshal([]byte(lines[0]), &m)

	if _, ok := m["timestamp"]; !ok {
		t.Error("cloned logger should preserve timestamp setting")
	}
}

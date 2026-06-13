package fastlog_test

import (
	"encoding/json"
	"testing"

	"github.com/nxtcoder17/fastlog"
)

func TestPackageLogger_DynamicLookup(t *testing.T) {
	// 1. Capture output from first logger
	cw1 := newCapture()
	l1 := fastlog.New().Writer(cw1).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(l1)

	fastlog.Info("message-one", "k1", "v1")

	// 2. Capture output from second logger
	cw2 := newCapture()
	l2 := fastlog.New().Writer(cw2).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(l2)

	fastlog.Info("message-two", "k2", "v2")

	// Verify l1 got message-one
	lines1 := cw1.Lines()
	if len(lines1) != 1 {
		t.Fatalf("expected 1 line in logger 1, got %d", len(lines1))
	}
	var m1 map[string]any
	if err := json.Unmarshal([]byte(lines1[0]), &m1); err != nil {
		t.Fatalf("failed to parse json for logger 1: %v", err)
	}
	if m1["message"] != "message-one" || m1["k1"] != "v1" {
		t.Errorf("unexpected content in logger 1: %v", m1)
	}

	// Verify l2 got message-two
	lines2 := cw2.Lines()
	if len(lines2) != 1 {
		t.Fatalf("expected 1 line in logger 2, got %d", len(lines2))
	}
	var m2 map[string]any
	if err := json.Unmarshal([]byte(lines2[0]), &m2); err != nil {
		t.Fatalf("failed to parse json for logger 2: %v", err)
	}
	if m2["message"] != "message-two" || m2["k2"] != "v2" {
		t.Errorf("unexpected content in logger 2: %v", m2)
	}
}

func TestPackageLogger_Methods(t *testing.T) {
	cw := newCapture()
	l := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).DebugMode(true).JSON()
	fastlog.SetDefaultLogger(l)

	fastlog.Debug("dbg")
	fastlog.Warn("wrn")
	fastlog.Error("err")

	lines := cw.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	// Test With
	cwWith := newCapture()
	lWith := fastlog.New().Writer(cwWith).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(lWith)
	wLogger := fastlog.With("added_key", "added_val")
	wLogger.Info("with_msg")

	linesWith := cwWith.Lines()
	if len(linesWith) != 1 {
		t.Fatalf("expected 1 line for with, got %d", len(linesWith))
	}
	var mWith map[string]any
	json.Unmarshal([]byte(linesWith[0]), &mWith)
	if mWith["added_key"] != "added_val" {
		t.Errorf("expected added_key in with output, got %v", mWith)
	}

	// Test Slog
	cwSlog := newCapture()
	lSlog := fastlog.New().Writer(cwSlog).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(lSlog)
	sLogger := fastlog.Slog()
	sLogger.Info("slog_msg")

	linesSlog := cwSlog.Lines()
	if len(linesSlog) != 1 {
		t.Fatalf("expected 1 line for slog, got %d", len(linesSlog))
	}

	// Test Clone
	cwClone := newCapture()
	lClone := fastlog.New().Writer(cwClone).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(lClone)
	builder := fastlog.Clone()
	clonedLogger := builder.JSON()
	clonedLogger.Info("clone_msg")

	linesClone := cwClone.Lines()
	if len(linesClone) != 1 {
		t.Fatalf("expected 1 line for clone, got %d", len(linesClone))
	}
}

func TestPackageLogger_NilGuard(t *testing.T) {
	cw := newCapture()
	l := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()
	fastlog.SetDefaultLogger(l)

	// Set nil should be safely ignored
	fastlog.SetDefaultLogger(nil)

	fastlog.Info("guard_msg")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	if m["message"] != "guard_msg" {
		t.Errorf("expected guard_msg, got %v", m["message"])
	}
}




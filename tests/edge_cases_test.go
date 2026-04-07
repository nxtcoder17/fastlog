package fastlog_test

import (
	"encoding/json"
	"fmt"
	"github.com/nxtcoder17/fastlog"
	"strings"
	"testing"
)

func TestEdgeCases_NilValue(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "nil_key", nil)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with nil value: %v", err)
	}
}

func TestEdgeCases_Unicode(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "unicode", "こんにちは世界")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with unicode: %v", err)
	}
	if m["unicode"] != "こんにちは世界" {
		t.Errorf("expected unicode value preserved, got %v", m["unicode"])
	}
}

func TestEdgeCases_SpecialChars(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "special", `quotes "and" backslashes \ and	newlines
here`)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with special chars: %v\noutput: %s", err, lines[0])
	}
}

func TestEdgeCases_EmptyAttrs(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with no attrs: %v", err)
	}
}

func TestEdgeCases_UnbalancedKV(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("panicked on unbalanced KV (acceptable): %v", r)
		}
	}()

	logger.Info("msg", "only_key")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestEdgeCases_EmptyMessage(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with empty message: %v", err)
	}
}

func TestEdgeCases_LargeAttrs(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	kv := make([]any, 0, 100)
	for i := 0; i < 50; i++ {
		kv = append(kv, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	logger.Info("msg", kv...)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with large attrs: %v", err)
	}
	if len(m) != 52 {
		t.Errorf("expected 52 fields, got %d", len(m))
	}
}

func TestEdgeCases_NestedMap(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "nested", map[string]any{
		"a": 1,
		"b": "two",
	})

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with nested map: %v", err)
	}
}

func TestEdgeCases_NestedSlice(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "slice", []any{1, "two", true})

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with nested slice: %v", err)
	}
}

func TestEdgeCases_ErrorType(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "err", fmt.Errorf("something went wrong"))

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON with error type: %v", err)
	}
	if !strings.Contains(m["err"].(string), "something went wrong") {
		t.Errorf("expected error message in output, got %v", m["err"])
	}
}

func TestEdgeCases_BoolAttrs(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "true_val", true, "false_val", false)

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["true_val"] != true {
		t.Errorf("expected true_val=true, got %v", m["true_val"])
	}
	if m["false_val"] != false {
		t.Errorf("expected false_val=false, got %v", m["false_val"])
	}
}

func TestEdgeCases_FloatAttrs(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).JSON()

	logger.Info("msg", "float64", 3.14159, "float32", float32(2.71))

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

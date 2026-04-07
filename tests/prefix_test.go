package fastlog_test

import (
	"github.com/nxtcoder17/fastlog"
	"strings"
	"testing"
)

func TestPrefix(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).Prefix("myprefix").JSON().Slog()

	logger.Info("msg", "key", "value")

	lines := cw.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "myprefix.key") {
		t.Errorf("expected prefixed key in: %s", lines[0])
	}
}

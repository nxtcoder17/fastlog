package fastlog_test

import (
	"encoding/json"
	"github.com/nxtcoder17/fastlog"
	"strings"
	"testing"
	"time"
)

func TestTimestamp_AllLoggers(t *testing.T) {
	formats := []struct {
		name   string
		format string
	}{
		{"RFC3339", time.RFC3339},
		{"RFC3339Nano", time.RFC3339Nano},
		{"DateTime", time.DateTime},
		{"DateOnly", time.DateOnly},
		{"TimeOnly", time.TimeOnly},
		{"RFC1123", time.RFC1123},
		{"RFC1123Z", time.RFC1123Z},
		{"RFC822", time.RFC822},
		{"RFC822Z", time.RFC822Z},
		{"Kitchen", time.Kitchen},
		{"ANSIC", time.ANSIC},
		{"Custom_SlashDate", "2006/01/02 15:04:05"},
		{"Custom_DashDate", "02-01-2006 15:04:05"},
		{"Custom_MonthName", "Jan 2, 2006 15:04:05"},
		{"Custom_Millis", "2006-01-02T15:04:05.000Z"},
	}

	originalFormat := fastlog.TimestampFormat

	t.Run("JSON", func(t *testing.T) {
		for _, tf := range formats {
			t.Run(tf.name, func(t *testing.T) {
				fastlog.TimestampFormat = tf.format
				t.Cleanup(func() {
					fastlog.TimestampFormat = originalFormat
				})

				cw := newCapture()
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(true).Colors(false).JSON()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				var m map[string]any
				if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
					t.Fatalf("invalid JSON: %v", err)
				}

				ts, ok := m["timestamp"]
				if !ok {
					t.Fatal("expected timestamp field")
				}
				tsStr, ok := ts.(string)
				if !ok {
					t.Fatalf("expected timestamp to be string, got %T", ts)
				}

				assertValidTimestamp(t, tsStr, tf.format)
			})
		}
	})

	t.Run("Logfmt", func(t *testing.T) {
		for _, tf := range formats {
			t.Run(tf.name, func(t *testing.T) {
				fastlog.TimestampFormat = tf.format
				t.Cleanup(func() {
					fastlog.TimestampFormat = originalFormat
				})

				cw := newCapture()
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(true).Colors(false).Logfmt()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				tsStr := extractLogfmtTimestamp(lines[0])
				if tsStr == "" {
					t.Fatal("expected non-empty timestamp")
				}

				assertValidTimestamp(t, tsStr, tf.format)
			})
		}
	})

	t.Run("Console", func(t *testing.T) {
		skipFormats := map[string]bool{
			"DateTime":         true,
			"RFC1123":          true,
			"RFC1123Z":         true,
			"RFC822":           true,
			"RFC822Z":          true,
			"ANSIC":            true,
			"Custom_SlashDate": true,
			"Custom_DashDate":  true,
			"Custom_MonthName": true,
		}

		for _, tf := range formats {
			t.Run(tf.name, func(t *testing.T) {
				if skipFormats[tf.name] {
					t.Skip("console output is space-delimited, can't extract timestamps containing spaces")
				}
				fastlog.TimestampFormat = tf.format
				t.Cleanup(func() {
					fastlog.TimestampFormat = originalFormat
				})

				cw := newCapture()
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(true).Colors(false).Console()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				parts := strings.SplitN(lines[0], " ", 2)
				if len(parts) < 1 {
					t.Fatal("empty console output")
				}
				tsStr := parts[0]

				assertValidTimestamp(t, tsStr, tf.format)
			})
		}
	})
}

func assertValidTimestamp(t *testing.T, tsStr, format string) {
	t.Helper()

	if tsStr == "" {
		t.Fatal("expected non-empty timestamp")
	}

	parsed, err := time.Parse(format, tsStr)
	if err != nil {
		t.Errorf("failed to parse timestamp %q with format %q: %v", tsStr, format, err)
	}
	if parsed.IsZero() {
		t.Errorf("parsed timestamp is zero")
	}
}

func extractLogfmtTimestamp(line string) string {
	idx := strings.Index(line, "timestamp=")
	if idx == -1 {
		return ""
	}
	val := line[idx+len("timestamp="):]
	if strings.HasPrefix(val, "\"") {
		end := strings.Index(val[1:], "\"")
		if end == -1 {
			return ""
		}
		return val[1 : end+1]
	}
	end := strings.IndexAny(val, " \n")
	if end == -1 {
		return strings.TrimSpace(val)
	}
	return val[:end]
}

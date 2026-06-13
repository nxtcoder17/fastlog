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
		name             string
		format           string
		timestampEnabled bool
	}{
		{"Timestamp Disabled", time.RFC3339, false},
		{"RFC3339", time.RFC3339, true},
		{"RFC3339Nano", time.RFC3339Nano, true},
		{"DateTime", time.DateTime, true},
		{"DateOnly", time.DateOnly, true},
		{"TimeOnly", time.TimeOnly, true},
		{"RFC1123", time.RFC1123, true},
		{"RFC1123Z", time.RFC1123Z, true},
		{"RFC822", time.RFC822, true},
		{"RFC822Z", time.RFC822Z, true},
		{"Kitchen", time.Kitchen, true},
		{"ANSIC", time.ANSIC, true},
		{"Custom_SlashDate", "2006/01/02 15:04:05", true},
		{"Custom_DashDate", "02-01-2006 15:04:05", true},
		{"Custom_MonthName", "Jan 2, 2006 15:04:05", true},
		{"Custom_Millis", "2006-01-02T15:04:05.000Z", true},
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
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(tf.timestampEnabled).Colors(false).JSON()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				var m map[string]any
				if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
					t.Fatalf("invalid JSON: %v", err)
				}

				if !tf.timestampEnabled {
					_, ok := m["timestamp"]
					if ok {
						t.Fatal("never expected timestamp field, when timestamp option is disabled")
					}
					return
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
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(tf.timestampEnabled).Colors(false).Logfmt()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				tsStr := extractLogfmtTimestamp(lines[0])

				if !tf.timestampEnabled {
					if tsStr != "" {
						t.Fatal("never expected non-empty timestamp, when timestamp option is disabled")
					}
					return
				}

				if tsStr == "" {
					t.Fatal("expected non-empty timestamp")
				}

				assertValidTimestamp(t, tsStr, tf.format)
			})
		}
	})

	t.Run("Console", func(t *testing.T) {
		for _, tf := range formats {
			t.Run(tf.name, func(t *testing.T) {
				fastlog.TimestampFormat = tf.format
				t.Cleanup(func() {
					fastlog.TimestampFormat = originalFormat
				})

				cw := newCapture()
				logger := fastlog.New().Writer(cw).Caller(false).Timestamp(tf.timestampEnabled).Colors(false).Console()

				logger.Info("msg")

				lines := cw.Lines()
				if len(lines) != 1 {
					t.Fatalf("expected 1 line, got %d", len(lines))
				}

				parts := strings.Split(lines[0], " | ")
				if len(parts) < 2 {
					t.Fatalf("unexpected console output format: %q", lines[0])
				}

				if !tf.timestampEnabled {
					if parts[0] != "INFO" {
						t.Fatalf("never expected timestamp when timestamp option is disabled, but got: %q", parts[0])
					}
					return
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

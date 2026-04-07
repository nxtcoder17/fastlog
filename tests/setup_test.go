package fastlog_test

import (
	"bytes"
	"strings"
)

// --- Helper: capture output from logger ---
type captureWriter struct {
	buf bytes.Buffer
}

func (w *captureWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *captureWriter) String() string {
	return w.buf.String()
}

func (w *captureWriter) Lines() []string {
	s := strings.TrimSpace(w.buf.String())
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func newCapture() *captureWriter {
	return &captureWriter{}
}

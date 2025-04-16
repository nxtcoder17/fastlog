package fastlog

import (
	"io"
	"sync"
)

type syncWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

// Write implements io.Writer.
func (sw *syncWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.writer.Write(p)
}

var _ io.Writer = (*syncWriter)(nil)

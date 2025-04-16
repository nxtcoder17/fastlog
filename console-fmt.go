package fastlog

import (
	"log/slog"
)

type consoleLogger struct {
	*loggerProps
}

// Debug implements loggerAPI.
func (c *consoleLogger) Debug(msg string, kv ...any) {
	c.handleLog(slog.LevelDebug, msg, kv...)
}

// Error implements loggerAPI.
func (c *consoleLogger) Error(msg string, kv ...any) {
	c.handleLog(slog.LevelError, msg, kv...)
}

// Info implements loggerAPI.
func (c *consoleLogger) Info(msg string, kv ...any) {
	c.handleLog(slog.LevelInfo, msg, kv...)
}

// Warn implements loggerAPI.
func (c *consoleLogger) Warn(msg string, kv ...any) {
	c.handleLog(slog.LevelWarn, msg, kv...)
}

// Slog implements loggerAPI.
func (c *consoleLogger) Slog() *slog.Logger {
	return slog.New(&consoleLoggerSlog{c.loggerProps})
}

var _ loggerAPI = (*consoleLogger)(nil)

func (l *consoleLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < l.LogLevel {
		return nil
	}

	buf := l.pool.Get()

	buf.AppendCaller(2)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()

	buf.Append('|')
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)

	if len(kv) >= 2 {
		buf.Append('\t')
	}

	for i := 1; i < len(kv); i += 2 {
		buf.AppendComponentSeparator()
		buf.AppendAttr(kv[i-1], kv[i])
	}

	buf.Append('\n')
	_, err := l.Writer.Write(buf.Bytes())
	l.pool.Put(buf)
	return err
}

package fastlog

import (
	"log/slog"
)

type consoleLogger struct {
	kv []any
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

// With implements loggerAPI.
func (c *consoleLogger) With(kv ...any) *Logger {
	return &Logger{
		loggerAPI: &consoleLogger{
			kv: kv,
			loggerProps: &loggerProps{
				attrs:   c.loggerProps.attrs,
				prefix:  c.loggerProps.prefix,
				pool:    NewPool(&c.loggerProps.Options),
				Options: c.loggerProps.Options,
			},
		},
	}
}

// Slog implements loggerAPI.
func (c *consoleLogger) Slog() *slog.Logger {
	return slog.New(&consoleLoggerSlog{c.loggerProps})
}

var _ loggerAPI = (*consoleLogger)(nil)

func (c *consoleLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < c.LogLevel {
		return nil
	}

	buf := c.pool.Get()

	if buf.AppendTimestamp() {
		buf.AppendComponentSeparator()
		buf.Append('|')
		buf.AppendComponentSeparator()
	}

	if buf.AppendCaller(2) {
		buf.AppendComponentSeparator()
		buf.Append('|')
		buf.AppendComponentSeparator()
	}

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()
	buf.Append('|')
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)

	if len(kv) >= 2 || len(c.kv) >= 2 {
		buf.Append('\t')
	}

	for i := 1; i < len(c.kv); i += 2 {
		buf.AppendComponentSeparator()
		buf.AppendAttr(c.kv[i-1], c.kv[i])
	}

	for i := 1; i < len(kv); i += 2 {
		buf.AppendComponentSeparator()
		buf.AppendAttr(kv[i-1], kv[i])
	}

	buf.Append('\n')
	_, err := c.Writer.Write(buf.Bytes())
	c.pool.Put(buf)
	return err
}

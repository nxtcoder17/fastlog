package fastlog

import (
	"log/slog"
)

type consoleLogger struct {
	kv []any
	prefix string
	opts *Options
	pool *Pool
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
func (c *consoleLogger) With(kv ...any) Logger {
	newKVs := make([]any, 0, len(kv) + len(c.kv))
	newKVs = append(newKVs, c.kv...)
	newKVs = append(newKVs, kv...)

	return &consoleLogger{
		kv: newKVs,
		prefix: c.prefix,
		opts: c.opts,
		pool: c.pool,
	}
}

// Slog implements loggerAPI.
func (c *consoleLogger) Slog() *slog.Logger {
	kv := make([]slog.Attr, 0, len(c.kv))
	for i := 1; i < len(c.kv); i++ {
		kv = append(kv, slog.Any(c.kv[i-1].(string), c.kv[i]))
	}

	return slog.New(&consoleLoggerSlog{
		kv: kv,
		prefix: c.prefix,
		pool: c.pool,
		opts: c.opts,
	})
}

// Clone implements loggerAPI.
func (c *consoleLogger) Clone() *loggerBuilder {
	return &loggerBuilder{
		prefix: c.prefix,
		options: c.opts,
	}
}

var _ Logger = (*consoleLogger)(nil)

func (c *consoleLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < c.opts.LogLevel {
		return nil
	}

	buf := c.pool.Get()

	if buf.AppendTimestamp() {
		buf.AppendComponentSeparator()
		buf.Append('|')
		buf.AppendComponentSeparator()
	}

	if buf.AppendCaller(2 + c.opts.SkipCallerFrames) {
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
	_, err := c.opts.Writer.Write(buf.Bytes())
	c.pool.Put(buf)
	return err
}

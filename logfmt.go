package fastlog

import (
	"log/slog"
)

type logfmtLogger struct {
	kv     []any
	prefix string
	opts   *Options
	pool   *Pool
}

// Debug implements LoggerAPI.
func (l *logfmtLogger) Debug(msg string, kv ...any) {
	l.handleLog(slog.LevelDebug, msg, kv...)
}

// Error implements LoggerAPI.
func (l *logfmtLogger) Error(msg string, kv ...any) {
	l.handleLog(slog.LevelError, msg, kv...)
}

// Info implements LoggerAPI.
func (l *logfmtLogger) Info(msg string, kv ...any) {
	l.handleLog(slog.LevelInfo, msg, kv...)
}

// Warn implements LoggerAPI.
func (l *logfmtLogger) Warn(msg string, kv ...any) {
	l.handleLog(slog.LevelWarn, msg, kv...)
}

// With implements loggerAPI.
func (l *logfmtLogger) With(kv ...any) Logger {
	attrs := make([]any, 0, len(kv)+len(l.kv))
	attrs = append(attrs, l.kv...)
	attrs = append(attrs, kv...)

	return &logfmtLogger{
		kv:   attrs,
		pool: l.pool,
		opts: l.opts,
	}
}

func (l *logfmtLogger) Slog() *slog.Logger {
	kv := make([]slog.Attr, 0, len(l.kv))
	for i := 1; i < len(l.kv); i += 2 {
		kv = append(kv, slog.Any(l.kv[i-1].(string), l.kv[i]))
	}

	return slog.New(&logfmtSlog{
		kv:   kv,
		pool: l.pool,
		opts: l.opts,
	})
}

// Clone implements loggerAPI.
func (l *logfmtLogger) Clone() *loggerBuilder {
	optsCopy := *l.opts
	return &loggerBuilder{
		options: &optsCopy,
		prefix:  l.prefix,
		kv:      l.kv,
	}
}

var _ Logger = (*logfmtLogger)(nil)

func (l *logfmtLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < l.opts.LogLevel {
		return nil
	}

	buf := l.pool.Get()

	needSep := false
	if buf.AppendTimestamp() {
		needSep = true
	}
	if buf.AppendCaller(2 + l.opts.SkipCallerFrames) {
		needSep = true
	}
	if needSep {
		buf.AppendComponentSeparator()
	}

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)
	buf.AppendComponentSeparator()

	for i := 1; i < len(l.kv); i += 2 {
		buf.AppendAttr(l.kv[i-1], l.kv[i])
		buf.AppendComponentSeparator()
	}

	for i := 1; i < len(kv); i += 2 {
		buf.AppendAttr(kv[i-1], kv[i])
		buf.AppendComponentSeparator()
	}

	buf.Append('\n')
	if _, err := l.opts.Writer.Write(buf.Bytes()); err != nil {
		return err
	}

	l.pool.Put(buf)
	return nil
}

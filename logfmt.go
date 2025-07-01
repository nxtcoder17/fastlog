package fastlog

import (
	"log/slog"
)

type logfmtLogger struct {
	kv []any
	*loggerProps
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
func (l *logfmtLogger) With(kv ...any) *Logger {
	return &Logger{
		loggerAPI: &logfmtLogger{
			kv: kv,
			loggerProps: &loggerProps{
				attrs:   l.loggerProps.attrs,
				prefix:  l.loggerProps.prefix,
				pool:    NewPool(&l.loggerProps.Options),
				Options: l.loggerProps.Options,
			},
		},
	}
}

func (l *logfmtLogger) Slog() *slog.Logger {
	return slog.New(&logfmtSlog{l.loggerProps})
}

// Clone implements loggerAPI.
func (c *logfmtLogger) Clone(options ...Options) *Logger {
	opts := Options{}
	if len(options) >= 1 {
		opts = options[0]
	}

	opts = c.loggerProps.Options.clone(opts)

	opts.withDefaults()

	return &Logger{
		loggerAPI: &logfmtLogger{
			kv: c.kv,
			loggerProps: &loggerProps{
				attrs:   c.loggerProps.attrs,
				prefix:  c.loggerProps.prefix,
				pool:    NewPool(&opts),
				Options: opts,
			},
		},
	}
}

var _ loggerAPI = (*logfmtLogger)(nil)

func (l *logfmtLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < l.LogLevel {
		return nil
	}

	buf := l.pool.Get()

	buf.AppendCaller(2 + l.SkipCallerFrames)
	buf.AppendComponentSeparator()

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
	_, err := l.Writer.Write(buf.Bytes())
	l.pool.Put(buf)
	return err
}

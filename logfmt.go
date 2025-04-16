package fastlog

import (
	"log/slog"
)

type logfmtLogger struct {
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

var _ loggerAPI = (*logfmtLogger)(nil)

func (l *logfmtLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < l.LogLevel {
		return nil
	}

	buf := l.pool.Get()

	buf.AppendCaller(2)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)
	buf.AppendComponentSeparator()

	for i := 1; i < len(kv); i += 2 {
		buf.AppendAttr(kv[i-1], kv[i])
		buf.AppendComponentSeparator()
	}

	buf.Append('\n')
	_, err := l.Writer.Write(buf.Bytes())
	l.pool.Put(buf)
	return err
}

func (l *logfmtLogger) Slog() *slog.Logger {
	return slog.New(&logfmtSlog{l.loggerProps})
}

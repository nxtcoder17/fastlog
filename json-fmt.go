package fastlog

import (
	"log/slog"
)

type jsonLogger struct {
	*loggerProps
}

// Debug implements loggerAPI.
func (j *jsonLogger) Debug(msg string, kv ...any) {
	j.handleLog(slog.LevelDebug, msg, kv...)
}

// Error implements loggerAPI.
func (j *jsonLogger) Error(msg string, kv ...any) {
	j.handleLog(slog.LevelError, msg, kv...)
}

// Info implements loggerAPI.
func (j *jsonLogger) Info(msg string, kv ...any) {
	j.handleLog(slog.LevelInfo, msg, kv...)
}

// Warn implements loggerAPI.
func (j *jsonLogger) Warn(msg string, kv ...any) {
	j.handleLog(slog.LevelWarn, msg, kv...)
}

// Slog implements loggerAPI.
func (j *jsonLogger) Slog() *slog.Logger {
	return slog.New(&jsonLoggerSlog{j.loggerProps})
}

var _ loggerAPI = (*jsonLogger)(nil)

func (j *jsonLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < j.LogLevel {
		return nil
	}

	buf := j.pool.Get()

	buf.Append("{")

	buf.AppendCaller(2)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)

	for i := 1; i < len(kv); i += 2 {
		buf.AppendComponentSeparator()
		buf.AppendAttr(kv[i-1], kv[i])
	}

	buf.Append('}')
	buf.Append('\n')

	_, err := j.Writer.Write(buf.Bytes())
	j.pool.Put(buf)
	return err
}

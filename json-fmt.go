package fastlog

import (
	"log/slog"
)

type jsonLogger struct {
	kv []any
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

// With implements loggerAPI.
func (j *jsonLogger) With(kv ...any) *Logger {
	return &Logger{
		loggerAPI: &jsonLogger{
			kv: kv,
			loggerProps: &loggerProps{
				attrs:   j.loggerProps.attrs,
				prefix:  j.loggerProps.prefix,
				pool:    j.loggerProps.pool,
				Options: j.loggerProps.Options,
			},
		},
	}
}

// Slog implements loggerAPI.
func (j *jsonLogger) Slog() *slog.Logger {
	return slog.New(&jsonLoggerSlog{j.loggerProps})
}

// Clone implements loggerAPI.
func (c *jsonLogger) Clone(options ...OptionFn) *Logger {
	opts := c.loggerProps.Options.clone(options...)

	return &Logger{
		loggerAPI: &jsonLogger{
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

var _ loggerAPI = (*jsonLogger)(nil)

func (j *jsonLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < j.LogLevel {
		return nil
	}

	buf := j.pool.Get()

	buf.Append("{")

	buf.AppendCaller(2 + j.SkipCallerFrames)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(msg)

	for i := 1; i < len(j.kv); i += 2 {
		buf.AppendComponentSeparator()
		buf.AppendAttr(j.kv[i-1], j.kv[i])
	}

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

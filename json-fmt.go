package fastlog

import (
	"log/slog"
)

type jsonLogger struct {
	kv []any
	prefix string
	opts *Options
	pool *Pool
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
func (j *jsonLogger) With(kv ...any) Logger {
	newKVs := make([]any, 0, len(kv) + len(j.kv))
	newKVs = append(newKVs, j.kv...)
	newKVs = append(newKVs, kv...)

	return &jsonLogger{
		kv: newKVs,
		prefix:  j.prefix,
		pool:    j.pool,
		opts: j.opts,
	}
}

// Slog implements loggerAPI.
func (j *jsonLogger) Slog() *slog.Logger {
	kv := make([]slog.Attr, 0, len(j.kv))
	for i := 1; i < len(j.kv); i++ {
		kv = append(kv, slog.Any(j.kv[i-1].(string), j.kv[i]))
	}

	return slog.New(&jsonLoggerSlog{
		kv: kv,
		prefix: j.prefix,
		pool: j.pool,
		opts: j.opts,
	})
}

// Clone implements loggerAPI.
func (j *jsonLogger) Clone() *loggerBuilder {
	return &loggerBuilder{
		options: j.opts,
		prefix: j.prefix,
	}
}

var _ Logger = (*jsonLogger)(nil)

func (j *jsonLogger) handleLog(level slog.Level, msg string, kv ...any) error {
	if level < j.opts.LogLevel {
		return nil
	}

	buf := j.pool.Get()

	buf.Append("{")

	buf.AppendCaller(2 + j.opts.SkipCallerFrames)
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

	_, err := j.opts.Writer.Write(buf.Bytes())
	j.pool.Put(buf)
	return err
}

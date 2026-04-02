package fastlog

import (
	"context"
	"log/slog"
)

type logfmtSlog struct {
	kv  []slog.Attr
	pool   *Pool
	prefix string
	opts *Options
}

var _ slog.Handler = (*logfmtSlog)(nil)

func (l *logfmtSlog) parseAttr(buf *Buffer, attr slog.Attr) {
	buf.Append(ColorKey)
	if l.prefix != "" {
		buf.Append(l.prefix)
		buf.Append(".")
	}
	buf.Append(attr.Key)
	buf.Append('=')
	buf.Append(ColorReset)

	switch attr.Value.Kind() {
	case slog.KindGroup:
		{
			if attr.Key == "" {
				for _, value := range attr.Value.Group() {
					if value.Value.Kind() != slog.KindGroup {
						l.parseAttr(buf, value)
					}
				}
				return
			}

			for _, value := range attr.Value.Group() {
				value.Key = attr.Key + "." + value.Key
				l.parseAttr(buf, value)
			}
		}
	case slog.KindString:
		{
			buf.AppendWithQuote(attr.Value.String())
		}
	default:
		{
			buf.Append(attr.Value.Any(), true)
		}
	}
}

// Enabled implements slog.Handler.
func (l *logfmtSlog) Enabled(ctx context.Context, lvl slog.Level) bool {
	return lvl >= l.opts.LogLevel
}

// Handle implements slog.Handler.
func (l *logfmtSlog) Handle(ctx context.Context, record slog.Record) error {
	buf := l.pool.Get()

	buf.AppendCaller(3 + l.opts.SkipCallerFrames)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(record.Level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(record.Message)
	buf.AppendComponentSeparator()

	c := 0
	record.AddAttrs(l.kv...)
	record.Attrs(func(a slog.Attr) bool {
		if c <= record.NumAttrs() {
			buf.AppendComponentSeparator()
		}
		c += 1
		l.parseAttr(buf, a)
		return true
	})

	buf.Append('\n')

	_, err := l.opts.Writer.Write(buf.Bytes())
	l.pool.Put(buf)
	return err
}

// WithAttrs implements slog.Handler.
func (l *logfmtSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	kv := make([]slog.Attr, 0, len(l.kv)+ len(attrs))
	kv = append(kv, l.kv...)
	kv = append(kv, attrs...)

	return &logfmtSlog{
		kv:    kv,
		prefix: l.prefix,
		pool: l.pool,
		opts: l.opts,
	}
}

// WithGroup implements slog.Handler.
func (l *logfmtSlog) WithGroup(name string) slog.Handler {
	return &logfmtSlog{
		kv:  l.kv,
		prefix: name + "." + l.prefix,
		pool:   l.pool,
		opts: l.opts,
	}
}

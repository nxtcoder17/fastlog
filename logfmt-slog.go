package fastlog

import (
	"context"
	"log/slog"
)

type logfmtSlog struct {
	*loggerProps
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
	return lvl >= l.LogLevel
}

// Handle implements slog.Handler.
func (l *logfmtSlog) Handle(ctx context.Context, record slog.Record) error {
	buf := l.pool.Get()

	buf.AppendCaller(3)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(record.Level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(record.Message)
	buf.AppendComponentSeparator()

	c := 0
	record.AddAttrs(l.attrs...)
	record.Attrs(func(a slog.Attr) bool {
		if c <= record.NumAttrs() {
			buf.AppendComponentSeparator()
		}
		c += 1
		l.parseAttr(buf, a)
		return true
	})

	buf.Append('\n')

	_, err := l.Writer.Write(buf.Bytes())
	l.pool.Put(buf)
	return err
}

// WithAttrs implements slog.Handler.
func (l *logfmtSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logfmtSlog{
		loggerProps: &loggerProps{
			pool:    NewPool(&l.Options),
			attrs:   append(l.attrs, attrs...),
			Options: l.Options,
		},
	}
}

// WithGroup implements slog.Handler.
func (l *logfmtSlog) WithGroup(name string) slog.Handler {
	return &logfmtSlog{
		loggerProps: &loggerProps{
			pool:   NewPool(&l.Options),
			attrs:  l.attrs,
			prefix: name + "." + l.prefix,

			Options: l.Options,
		},
	}
}

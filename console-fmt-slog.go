package fastlog

import (
	"context"
	"log/slog"
)

type consoleLoggerSlog struct {
	*loggerProps
}

func (l *consoleLoggerSlog) parseAttr(buf *Buffer, attr slog.Attr) {
	buf.AppendAttrKeyColor()
	if l.prefix != "" {
		buf.Append(l.prefix)
		buf.Append(".")
	}
	buf.Append(attr.Key)
	buf.AppendAttrKeyColorReset()

	buf.AppendAttrSeparator()

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
func (l *consoleLoggerSlog) Enabled(ctx context.Context, lvl slog.Level) bool {
	return lvl >= l.LogLevel
}

// Handle implements slog.Handler.
func (l *consoleLoggerSlog) Handle(ctx context.Context, record slog.Record) error {
	buf := l.pool.Get()

	buf.AppendCaller(3 + l.SkipCallerFrames)
	buf.AppendComponentSeparator()
	buf.AppendLogLevel(record.Level)
	buf.AppendComponentSeparator()
	buf.Append('|')
	buf.AppendComponentSeparator()

	buf.AppendMsg(record.Message)

	if record.NumAttrs() >= 2 {
		buf.Append('\t')
	}

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
func (l *consoleLoggerSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &consoleLoggerSlog{
		loggerProps: &loggerProps{
			pool:    NewPool(&l.Options),
			attrs:   append(l.attrs, attrs...),
			Options: l.Options,
		},
	}
}

// WithGroup implements slog.Handler.
func (l *consoleLoggerSlog) WithGroup(name string) slog.Handler {
	return &consoleLoggerSlog{
		loggerProps: &loggerProps{
			pool:   NewPool(&l.Options),
			attrs:  l.attrs,
			prefix: name + "." + l.prefix,

			Options: l.Options,
		},
	}
}

var _ slog.Handler = (*consoleLoggerSlog)(nil)

package fastlog

import (
	"context"
	"log/slog"
)

type jsonLoggerSlog struct {
	*loggerProps
}

func (j *jsonLoggerSlog) parseAttr(buf *Buffer, attr slog.Attr) {
	if j.prefix != "" {
		buf.AppendAttrKey(j.prefix + "." + attr.Key)
	} else {
		buf.AppendAttrKey(attr.Key)
	}
	buf.AppendAttrSeparator()

	switch attr.Value.Kind() {
	case slog.KindGroup:
		{
			if attr.Key == "" {
				for _, value := range attr.Value.Group() {
					if value.Value.Kind() != slog.KindGroup {
						j.parseAttr(buf, value)
					}
				}
				return
			}

			for _, value := range attr.Value.Group() {
				value.Key = attr.Key + "." + value.Key
				j.parseAttr(buf, value)
			}
		}
	case slog.KindString:
		{
			buf.AppendWithQuote(attr.Value.String())
		}
	default:
		{
			buf.Append(attr.Value.Any())
		}
	}
}

// Enabled implements slog.Handler.
func (j *jsonLoggerSlog) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= j.LogLevel
}

// Handle implements slog.Handler.
func (j *jsonLoggerSlog) Handle(_ context.Context, record slog.Record) error {
	buf := j.pool.Get()
	buf.Append("{")
	buf.AppendCaller(3)
	buf.AppendComponentSeparator()

	buf.AppendLogLevel(record.Level)
	buf.AppendComponentSeparator()

	buf.AppendMsg(record.Message)

	c := 0
	record.AddAttrs(j.attrs...)
	record.Attrs(func(a slog.Attr) bool {
		if c <= record.NumAttrs() {
			buf.AppendComponentSeparator()
		}
		c += 1
		j.parseAttr(buf, a)
		return true
	})

	buf.Append('}')
	buf.Append('\n')
	_, err := j.Writer.Write(buf.Bytes())
	j.pool.Put(buf)
	return err
}

// WithAttrs implements slog.Handler.
func (j *jsonLoggerSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &jsonLoggerSlog{
		loggerProps: &loggerProps{
			pool:    &Pool{},
			attrs:   append(j.attrs, attrs...),
			Options: j.Options,
		},
	}
}

// WithGroup implements slog.Handler.
func (j *jsonLoggerSlog) WithGroup(name string) slog.Handler {
	return &jsonLoggerSlog{
		loggerProps: &loggerProps{
			pool:    &Pool{},
			attrs:   j.attrs,
			prefix:  name + "." + j.prefix,
			Options: j.Options,
		},
	}
}

var _ slog.Handler = (*jsonLoggerSlog)(nil)

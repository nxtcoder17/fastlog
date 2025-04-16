package fastlog

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Pool struct {
	sync.Pool
}

func NewPool(opts *Options) *Pool {
	return &Pool{
		Pool: sync.Pool{
			New: func() any {
				return &Buffer{
					store:   make([]byte, 0, 256),
					Options: opts,
				}
			},
		},
	}
}

func (p *Pool) Get() *Buffer {
	if v := p.Pool.Get(); v != nil {
		return v.(*Buffer)
	}
	return &Buffer{}
}

func (p *Pool) Put(msg *Buffer) {
	msg.Reset()
	p.Pool.Put(msg)
}

type Buffer struct {
	store []byte
	*Options
}

func (buf *Buffer) Appendf(s string, args ...any) {
	buf.store = fmt.Appendf(buf.store, s, args...)
}

func (buf *Buffer) AppendLogLevel(lvl slog.Level) {
	switch buf.Format {
	case ConsoleFormat:
	default:
		buf.AppendAttrKey(buf.LevelFieldKey)
		buf.AppendAttrSeparator()
	}

	switch lvl {
	case slog.LevelDebug:
		if buf.EnableColors {
			buf.store = append(buf.store, FgWhite...)
		}
		buf.AppendAttrValue("DEBUG")

	case slog.LevelInfo:
		if buf.EnableColors {
			buf.store = append(buf.store, FgGreen...)
		}
		buf.AppendAttrValue("INFO")

	case slog.LevelWarn:
		if buf.EnableColors {
			buf.store = append(buf.store, FgYellow...)
		}
		buf.AppendAttrValue("WARN")
	case slog.LevelError:
		if buf.EnableColors {
			buf.store = append(buf.store, FgRed...)
		}
		buf.AppendAttrValue("ERROR")
	}

	if buf.EnableColors {
		buf.store = append(buf.store, ColorReset...)
	}
}

func (buf *Buffer) AppendAttrSeparator() {
	if buf.EnableColors {
		buf.store = append(buf.store, ColorSeparator...)
	}
	switch buf.Format {
	case JSONFormat:
		buf.store = append(buf.store, ':')
	default:
		buf.store = append(buf.store, '=')
	}
	if buf.EnableColors {
		buf.store = append(buf.store, ColorReset...)
	}
}

func (buf *Buffer) AppendComponentSeparator() {
	switch buf.Format {
	case JSONFormat:
		buf.store = append(buf.store, ',')
	default:
		buf.store = append(buf.store, ' ')
	}
}

func (buf *Buffer) AppendMsg(msg string) {
	switch buf.Format {
	case ConsoleFormat:
	case JSONFormat, LogfmtFormat:
		buf.AppendAttrKey(buf.MessageFieldKey)
		buf.AppendAttrSeparator()
	}

	if buf.EnableColors {
		buf.store = append(buf.store, ColorMessage...)
	}
	buf.AppendAttrValue(msg)
	if buf.EnableColors {
		buf.store = append(buf.store, ColorReset...)
	}
}

func (buf *Buffer) AppendAttrKey(key any) {
	if buf.EnableColors {
		buf.store = append(buf.store, ColorKey...)
	}

	switch buf.Format {
	case JSONFormat:
		buf.AppendWithQuote(key)
	default:
		buf.Append(key)
	}

	if buf.EnableColors {
		buf.store = append(buf.store, ColorReset...)
	}
}

func (buf *Buffer) AppendAttrKeyColor() {
	if buf.EnableColors {
		buf.store = append(buf.store, ColorKey...)
	}
}

func (buf *Buffer) AppendAttrKeyColorReset() {
	if buf.EnableColors {
		buf.store = append(buf.store, ColorReset...)
	}
}

func (buf *Buffer) AppendAttrValue(value any) {
	switch buf.Format {
	case JSONFormat:
		buf.AppendWithQuote(value)
	default:
		buf.Append(value)
	}
}

func (buf *Buffer) AppendCaller(skip int) {
	if buf.ShowCaller {
		_, file, line, ok := runtime.Caller(skip + 1)
		if ok {
			switch buf.Format {
			case ConsoleFormat:
			default:
				buf.AppendAttrKey(buf.CallerFieldKey)
				buf.AppendAttrSeparator()
				buf.Append('"')
				defer buf.Append('"')
			}

			if buf.EnableColors {
				buf.Append(FgBrightBlack)
			}

			trimCallerPath(buf, file, 2)
			buf.Append(':')
			buf.Append(line)
		}
	}
}

func (buf *Buffer) AppendAttr(key, value any) {
	buf.AppendAttrKey(key)
	buf.AppendAttrSeparator()
	buf.AppendAttrValue(value)
}

func (buf *Buffer) Append(b any, noQuote ...bool) {
	switch v := b.(type) {
	case byte:
		buf.store = append(buf.store, v)
	case []byte:
		buf.store = append(buf.store, v...)
	case string:
		buf.store = append(buf.store, v...)

	case int64:
		buf.store = strconv.AppendInt(buf.store, v, 10)
	case int:
		buf.store = strconv.AppendInt(buf.store, int64(v), 10)
	case int32:
		if v >= 0 && v <= 127 {
			buf.store = append(buf.store, byte(v))
		} else {
			buf.store = strconv.AppendInt(buf.store, int64(v), 10)
		}

	case float64:
		buf.store = strconv.AppendFloat(buf.store, v, 'f', -1, 64)
	case float32:
		buf.store = strconv.AppendFloat(buf.store, float64(v), 'f', -1, 32)

	case bool:
		buf.store = strconv.AppendBool(buf.store, v)

	case error:
		buf.store = appendStringWithQuotes(buf.store, []byte(v.Error()))

	case time.Time:
		buf.store = append(buf.store, v.Format(buf.TimestampFormat)...)

	case []any:
		appendSliceIntoBuf(buf, v)
	case []int:
		appendSliceIntoBuf(buf, v)
	case []int32:
		appendSliceIntoBuf(buf, v)
	case []int64:
		appendSliceIntoBuf(buf, v)
	case []bool:
		appendSliceIntoBuf(buf, v)
	case []float32:
		appendSliceIntoBuf(buf, v)
	case []float64:
		appendSliceIntoBuf(buf, v)
	case []string:
		appendSliceIntoBuf(buf, v)
	case []map[string]any:
		appendSliceIntoBuf(buf, v)

	case map[string]any:
		appendMapIntoBuf(buf, v)

	case map[string]string:
		appendMapIntoBuf(buf, v)

	case map[string]int:
		appendMapIntoBuf(buf, v)

	case map[string]bool:
		appendMapIntoBuf(buf, v)

	case map[string]float64:
		appendMapIntoBuf(buf, v)

	default:
		if len(noQuote) > 0 && noQuote[0] {
			buf.store = fmt.Appendf(buf.store, "%#v", v)
			return
		}
		// INFO: this is just a hack to have something like fmt.Sprintf("") but not allocating any additional space, and using our pool only to store it, and later override it
		l := len(buf.store)
		buf.store = fmt.Appendf(buf.store, "%#v", v)
		l1 := len(buf.store)
		m := buf.store[l:]

		buf.store = appendStringWithQuotes(buf.store, m)
		l2 := len(buf.store)

		diff := l2 - l1

		for i := 0; i < diff; i++ {
			buf.store[l+i] = buf.store[l1+i]
		}

		buf.store = buf.store[:l+diff]
	}
}

func appendMapIntoBuf[K comparable, T any](buf *Buffer, m map[K]T) {
	buf.store = append(buf.store, '{')
	addComma := false
	for k, v := range m {
		if addComma {
			buf.store = append(buf.store, ',')
		}
		buf.AppendWithQuote(k)
		// buf.Append(k)
		buf.store = append(buf.store, ':')
		buf.AppendWithQuote(v)
		// buf.Append(v)

		addComma = true
	}
	buf.store = append(buf.store, '}')
}

func appendSliceIntoBuf[T any](buf *Buffer, arr []T) {
	buf.store = append(buf.store, '[')
	for i := range arr {
		if i > 0 {
			buf.store = append(buf.store, ',')
		}
		buf.AppendWithQuote(arr[i])
	}
	buf.store = append(buf.store, ']')
}

func (buf *Buffer) AppendWithQuote(v any) {
	switch val := v.(type) {
	case []byte:
		buf.store = appendStringWithQuotes(buf.store, val)
		// buf.store = strconv.AppendQuote(buf.store, string(val))
	case string:
		buf.store = appendStringWithQuotes(buf.store, []byte(val))
		// buf.store = strconv.AppendQuote(buf.store, val)
	default:
		buf.Append(v)
	}
}

func (buf *Buffer) Bytes() []byte {
	return buf.store
}

func (buf *Buffer) Reset() {
	buf.store = buf.store[:0]
}

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
	newBuffer func() *Buffer
	sync.Pool
}

func NewPool(opts *Options) *Pool {
	newBuffer := func() *Buffer {
		return &Buffer{
			store:   make([]byte, 0, InitialBufSize),
			Options: opts,
		}
	}
	return &Pool{
		newBuffer: newBuffer,
		Pool: sync.Pool{
			New: func() any {
				return newBuffer()
			},
		},
	}
}

func (p *Pool) Get() *Buffer {
	if v := p.Pool.Get(); v != nil {
		return v.(*Buffer)
	}
	return p.newBuffer()
}

func (p *Pool) Put(buf *Buffer) {
	buf.Reset()
	// If buffer has grown too large, replace it with a new smaller buffer
	if cap(buf.store) > MaxBufSize {
		buf = p.newBuffer()
	}
	p.Pool.Put(buf)
}

type Buffer struct {
	store []byte
	*Options
}

func (buf *Buffer) Appendf(s string, args ...any) {
	buf.store = fmt.Appendf(buf.store, s, args...)
}

func (buf *Buffer) AppendLogLevel(lvl slog.Level) bool {
	switch buf.Format {
	case logFormatConsole:
	default:
		buf.AppendAttrKey(LevelFieldKey)
		buf.AppendAttrSeparator()
	}

	switch lvl {
	case slog.LevelDebug:
		if !buf.EnableColors {
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

	return true
}

func (buf *Buffer) AppendAttrSeparator() {
	if buf.EnableColors {
		buf.store = append(buf.store, ColorSeparator...)
	}
	switch buf.Format {
	case logFormatJSON:
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
	case logFormatJSON:
		buf.store = append(buf.store, ',')
	default:
		buf.store = append(buf.store, ' ')
	}
}

func (buf *Buffer) AppendMsg(msg string) {
	switch buf.Format {
	case logFormatJSON, logFormatLogFmt:
		buf.AppendAttrKey(MessageFieldKey)
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
	case logFormatJSON:
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
	case logFormatJSON:
		buf.AppendWithQuote(value)
	default:
		buf.Append(value)
	}
}

func (buf *Buffer) AppendCaller(skip int) bool {
	if buf.ShowCaller {
		_, file, line, ok := runtime.Caller(skip + 1)
		if ok {
			switch buf.Format {
			case logFormatConsole:
			default:
				buf.AppendAttrKey(CallerFieldKey)
				buf.AppendAttrSeparator()
				buf.Append('"')
				defer buf.Append('"')
			}

			if buf.EnableColors {
				buf.Append(FgBrightBlack)
			}

			trimCallerPath(buf, file, 2)
			buf.Append(':')
			if buf.Format == logFormatConsole {
				buf.Appendf("%-3d", line)
			} else {
				buf.Append(line)
			}
		}
		return true
	}
	return false
}

func (buf *Buffer) AppendTimestamp() bool {
	if buf.ShowTimestamp {
		switch buf.Format {
		case logFormatConsole:
		default:
			buf.AppendAttrKey(TimestampFieldKey)
			buf.AppendAttrSeparator()
			buf.Append('"')
			defer buf.Append('"')
		}

		if buf.EnableColors {
			buf.Append(FgBrightBlack)
		}

		buf.Append(time.Now().Format(TimestampFormat))
		return true
	}

	return false
}

func (buf *Buffer) AppendAttr(key, value any) {
	buf.AppendAttrKey(key)
	buf.AppendAttrSeparator()
	buf.AppendAttrValue(value)
}

func (buf *Buffer) Append(b any, quote ...bool) {
	switch v := b.(type) {
	case byte:
		buf.store = append(buf.store, v)
	case []byte:
		buf.store = append(buf.store, v...)
	case string:
		buf.store = append(buf.store, v...)
	case time.Time:
		buf.store = append(buf.store, v.Format(TimestampFormat)...)
	case fmt.Stringer:
		buf.store = append(buf.store, v.String()...)

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
		if len(quote) > 0 && quote[0] {
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
			return
		}
		buf.store = fmt.Appendf(buf.store, "%#v", v)
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
	case string:
		buf.store = appendStringWithQuotes(buf.store, []byte(val))
	default:
		buf.Append(v, true)
	}
}

func (buf *Buffer) Bytes() []byte {
	return buf.store
}

func (buf *Buffer) Reset() {
	buf.store = buf.store[:0]
}

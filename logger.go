package fastlog

import (
	"io"
	"log/slog"
	"os"
	"time"
)

type loggerProps struct {
	attrs  []slog.Attr
	prefix string
	pool   *Pool

	Options
}

type logformat string

const (
	ConsoleFormat logformat = "console"
	JSONFormat    logformat = "json"
	LogfmtFormat  logformat = "logfmt"
)

type Options struct {
	Writer   io.Writer
	Format   logformat
	LogLevel slog.Level

	// ShowDebugLogs sets LogLevel to Debug
	ShowDebugLogs bool

	ShowCaller   bool
	EnableColors bool

	TimestampFieldKey string
	TimestampFormat   string

	MessageFieldKey string
	LevelFieldKey   string
	CallerFieldKey  string
}

func New(options ...Options) *Logger {
	opts := Options{}
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Writer == nil {
		opts.Writer = os.Stderr
	}

	opts.Writer = &syncWriter{writer: opts.Writer}

	if opts.TimestampFieldKey == "" {
		opts.TimestampFieldKey = "ts"
	}

	if opts.TimestampFormat == "" {
		opts.TimestampFormat = time.RFC850
	}

	if opts.MessageFieldKey == "" {
		opts.MessageFieldKey = "message"
	}

	if opts.LevelFieldKey == "" {
		opts.LevelFieldKey = "level"
	}

	if opts.CallerFieldKey == "" {
		opts.CallerFieldKey = "caller"
	}

	if opts.ShowDebugLogs {
		opts.LogLevel = slog.LevelDebug
	}

	props := &loggerProps{
		attrs:  nil,
		prefix: "",
		pool:   NewPool(&opts),

		Options: opts,
	}

	switch opts.Format {
	case ConsoleFormat:
		return &Logger{loggerAPI: &consoleLogger{props}}
	case LogfmtFormat:
		return &Logger{loggerAPI: &logfmtLogger{props}}
	case JSONFormat:
		return &Logger{loggerAPI: &jsonLogger{props}}
	default:
		return &Logger{loggerAPI: &consoleLogger{props}}
	}
}

type loggerAPI interface {
	Info(msg string, kv ...any)
	Debug(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)

	Slog() *slog.Logger
}

type Logger struct {
	loggerAPI
}

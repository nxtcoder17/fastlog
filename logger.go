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
	ShowCaller    bool
	ShowTimestamp bool

	SkipCallerFrames int

	EnableColors bool

	TimestampFieldKey string
	TimestampFormat   string

	MessageFieldKey string
	LevelFieldKey   string
	CallerFieldKey  string
}

func (opts *Options) withDefaults() {
	if opts.Writer == nil {
		opts.Writer = os.Stderr
	}

	opts.Writer = &syncWriter{writer: opts.Writer}

	if opts.TimestampFieldKey == "" {
		opts.TimestampFieldKey = "ts"
	}

	if opts.TimestampFormat == "" {
		opts.TimestampFormat = time.RFC3339
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
}

func (opts Options) clone(opt2 Options) Options {
	if opt2.Writer != nil {
		opts.Writer = opt2.Writer
	}

	if opt2.TimestampFieldKey != "" {
		opts.TimestampFieldKey = opt2.TimestampFieldKey
	}

	if opt2.TimestampFormat != "" {
		opts.TimestampFormat = opt2.TimestampFormat
	}

	if opt2.MessageFieldKey != "" {
		opts.MessageFieldKey = opt2.MessageFieldKey
	}

	if opt2.LevelFieldKey != "" {
		opts.LevelFieldKey = opt2.LevelFieldKey
	}

	if opt2.CallerFieldKey != "" {
		opts.CallerFieldKey = opt2.CallerFieldKey
	}

	if opts.ShowDebugLogs || opt2.ShowDebugLogs {
		opts.LogLevel = slog.LevelDebug
	}

	if opt2.SkipCallerFrames > 0 {
		opts.SkipCallerFrames = opt2.SkipCallerFrames
	}

	return opts
}

func New(options ...Options) *Logger {
	opts := Options{}
	if len(options) > 0 {
		opts = options[0]
	}

	opts.withDefaults()

	props := &loggerProps{
		attrs:  nil,
		prefix: "",
		pool:   NewPool(&opts),

		Options: opts,
	}

	switch opts.Format {
	case ConsoleFormat:
		return &Logger{loggerAPI: &consoleLogger{kv: nil, loggerProps: props}}
	case LogfmtFormat:
		return &Logger{loggerAPI: &logfmtLogger{kv: nil, loggerProps: props}}
	case JSONFormat:
		return &Logger{loggerAPI: &jsonLogger{kv: nil, loggerProps: props}}
	default:
		return &Logger{loggerAPI: &consoleLogger{kv: nil, loggerProps: props}}
	}
}

type loggerAPI interface {
	Info(msg string, kv ...any)
	Debug(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)

	With(kv ...any) *Logger
	Slog() *slog.Logger

	Clone(option ...Options) *Logger
}

type Logger struct {
	loggerAPI
}

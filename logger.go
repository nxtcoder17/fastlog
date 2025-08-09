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
	logFormatConsole logformat = "console"
	logFormatJSON    logformat = "json"
	logFormatLogFmt  logformat = "logfmt"
)

var (
	LevelFieldKey   string = "level"
	CallerFieldKey  string = "caller"
	MessageFieldKey string = "message"

	TimestampFieldKey string = "timestamp"
	TimestampFormat   string = time.RFC3339
)

var (
	InitialBufSize int = 256
	MaxBufSize     int = 4096
)

type Options struct {
	Writer io.Writer
	Format logformat

	LogLevel slog.Level

	// ShowDebugLogs and LogLevel are mutually exclusive, as ShowDebugLogs set LogLevel to Debug regardless of LogLevel field being set or not
	ShowDebugLogs bool

	ShowCaller bool

	ShowTimestamp bool
	EnableColors  bool

	SkipCallerFrames int
}

type OptionFn func(opt *Options)

func WithWriter(writer io.Writer) OptionFn {
	return func(opt *Options) {
		opt.Writer = writer
	}
}

func WithLogLevel(level slog.Level) OptionFn {
	return func(opt *Options) {
		opt.LogLevel = level
	}
}

func WithoutCaller() OptionFn {
	return func(opt *Options) {
		opt.ShowCaller = false
	}
}

func ShowDebugLogs(enabled bool) OptionFn {
	return func(opt *Options) {
		opt.ShowDebugLogs = enabled
	}
}

func SkipCallerFrames(count int) OptionFn {
	return func(opt *Options) {
		opt.SkipCallerFrames = count
	}
}

func Json() OptionFn {
	return func(opt *Options) {
		opt.Format = logFormatJSON
	}
}

func Logfmt() OptionFn {
	return func(opt *Options) {
		opt.Format = logFormatLogFmt
	}
}

func Console() OptionFn {
	return func(opt *Options) {
		opt.Format = logFormatConsole
	}
}

func WithoutColors() OptionFn {
	return func(opt *Options) {
		opt.EnableColors = false
	}
}

func WithoutTimestamp() OptionFn {
	return func(opt *Options) {
		opt.ShowTimestamp = false
	}
}

func (opts Options) clone(options ...OptionFn) Options {
	for _, optFn := range options {
		optFn(&opts)
	}

	if _, ok := opts.Writer.(*syncWriter); !ok {
		opts.Writer = &syncWriter{writer: opts.Writer}
	}

	return opts
}

func New(options ...OptionFn) *Logger {
	opts := Options{
		Writer:           os.Stderr,
		Format:           logFormatConsole,
		LogLevel:         slog.LevelInfo,
		ShowDebugLogs:    false,
		ShowCaller:       true,
		ShowTimestamp:    true,
		EnableColors:     true,
		SkipCallerFrames: 0,
	}

	for _, optFn := range options {
		optFn(&opts)
	}

	opts.Writer = &syncWriter{writer: opts.Writer}

	props := &loggerProps{
		attrs:  nil,
		prefix: "",
		pool:   NewPool(&opts),

		Options: opts,
	}

	switch opts.Format {
	case logFormatConsole:
		return &Logger{loggerAPI: &consoleLogger{kv: nil, loggerProps: props}}
	case logFormatLogFmt:
		return &Logger{loggerAPI: &logfmtLogger{kv: nil, loggerProps: props}}
	case logFormatJSON:
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

	Clone(option ...OptionFn) *Logger
}

type Logger struct {
	loggerAPI
}

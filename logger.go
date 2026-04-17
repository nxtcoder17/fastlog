package fastlog

import (
	"github.com/nxtcoder17/fastlog/types"
	"io"
	"log/slog"
	"os"
	"time"
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
	Format types.LogFormat

	// LogLevel defaults to Info
	LogLevel slog.Level

	WithCaller bool

	WithTimestamp bool
	WithColors    bool

	SkipCallerFrames int
	VerbosityLevel   int
}

func New() *loggerBuilder {
	return &loggerBuilder{
		options: &Options{
			Writer:         &syncWriter{writer: os.Stderr},
			Format:         types.LogFormatConsole,
			LogLevel:       slog.LevelInfo,
			WithCaller:     true,
			WithTimestamp:  true,
			WithColors:     true,
			VerbosityLevel: 1,
		},
	}
}

func (l *loggerBuilder) JSON() Logger {
	return &jsonLogger{
		kv:     nil,
		prefix: l.prefix,
		opts:   l.options,
		pool: newPool(&bufferPoolOptions{
			WithTimestamp: l.options.WithTimestamp,
			WithCaller:    l.options.WithCaller,
			WithColors:    l.options.WithColors,
			LogFormat:     types.LogFormatJSON,
		}),
	}
}

func (l *loggerBuilder) Logfmt() Logger {
	return &logfmtLogger{
		kv:   nil,
		opts: l.options,
		pool: newPool(&bufferPoolOptions{
			WithTimestamp: l.options.WithTimestamp,
			WithCaller:    l.options.WithCaller,
			WithColors:    l.options.WithColors,
			LogFormat:     types.LogFormatLogfmt,
		}),
	}
}

func (l *loggerBuilder) Console() Logger {
	return &consoleLogger{
		kv:   nil,
		opts: l.options,
		pool: newPool(&bufferPoolOptions{
			WithTimestamp: l.options.WithTimestamp,
			WithCaller:    l.options.WithCaller,
			WithColors:    l.options.WithColors,
			LogFormat:     types.LogFormatConsole,
		}),
	}
}

func (l *loggerBuilder) Writer(writer io.Writer) *loggerBuilder {
	l.options.Writer = &syncWriter{writer: writer}
	return l
}

func (l *loggerBuilder) Prefix(p string) *loggerBuilder {
	l.prefix = p
	return l
}

func (l *loggerBuilder) Timestamp(show bool) *loggerBuilder {
	l.options.WithTimestamp = show
	return l
}

func (l *loggerBuilder) Caller(show bool) *loggerBuilder {
	l.options.WithCaller = show
	return l
}

func (l *loggerBuilder) Colors(show bool) *loggerBuilder {
	l.options.WithColors = show
	return l
}

func (l *loggerBuilder) LogLevel(level slog.Level) *loggerBuilder {
	l.options.LogLevel = level
	return l
}

// DebugMode is an alias to LogLevel(slog.LevelDebug)
func (l *loggerBuilder) DebugMode(enable bool) *loggerBuilder {
	l.options.LogLevel = slog.LevelDebug
	return l
}

func (l *loggerBuilder) SkipCallerFrames(n int) *loggerBuilder {
	l.options.SkipCallerFrames = n
	return l
}

func (l *loggerBuilder) Verbosity(n int) *loggerBuilder {
	l.options.VerbosityLevel = n
	return l
}

type loggerBuilder struct {
	prefix  string
	options *Options
}

type Logger interface {
	Info(msg string, kv ...any)
	Debug(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)

	// Verbosity() returns the verbosity level of the logger
	// Verbosity() int

	With(kv ...any) Logger
	Slog() *slog.Logger

	Clone() *loggerBuilder
}

func JSON() Logger {
	return New().JSON()
}

func Logfmt() Logger {
	return New().Logfmt()
}

func Console() Logger {
	return New().Console()
}

var defaultLogger Logger

func SetDefault(logger Logger) {
	defaultLogger = logger
}

func Default() Logger {
	if defaultLogger != nil {
		return defaultLogger
	}

	return Console()
}

package fastlog

import (
	"log/slog"
)

var logger Logger = Console()

func SetDefaultLogger(l Logger) {
	if l != nil {
		logger = l
	}
}

func Debug(msg string, kv ...any) {
	logger.Debug(msg, kv...)
}

func Info(msg string, kv ...any) {
	logger.Info(msg, kv...)
}

func Warn(msg string, kv ...any) {
	logger.Warn(msg, kv...)
}

func Error(msg string, kv ...any) {
	logger.Error(msg, kv...)
}

func With(kv ...any) Logger {
	return logger.With(kv...)
}

func Slog() *slog.Logger {
	return logger.Slog()
}

func Clone() *loggerBuilder {
	return logger.Clone()
}

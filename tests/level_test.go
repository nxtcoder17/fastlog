package fastlog_test

import (
	"github.com/nxtcoder17/fastlog"
	"log/slog"
	"testing"
)

func TestLogLevel_Filtering(t *testing.T) {
	tests := []struct {
		name     string
		logLevel slog.Level
		call     func(fastlog.Logger)
		expect   bool
	}{
		{"Debug_at_DebugLevel", slog.LevelDebug, func(l fastlog.Logger) { l.Debug("msg") }, true},
		{"Info_at_DebugLevel", slog.LevelDebug, func(l fastlog.Logger) { l.Info("msg") }, true},
		{"Debug_at_InfoLevel", slog.LevelInfo, func(l fastlog.Logger) { l.Debug("msg") }, false},
		{"Info_at_InfoLevel", slog.LevelInfo, func(l fastlog.Logger) { l.Info("msg") }, true},
		{"Warn_at_InfoLevel", slog.LevelInfo, func(l fastlog.Logger) { l.Warn("msg") }, true},
		{"Debug_at_WarnLevel", slog.LevelWarn, func(l fastlog.Logger) { l.Debug("msg") }, false},
		{"Info_at_WarnLevel", slog.LevelWarn, func(l fastlog.Logger) { l.Info("msg") }, false},
		{"Warn_at_WarnLevel", slog.LevelWarn, func(l fastlog.Logger) { l.Warn("msg") }, true},
		{"Debug_at_ErrorLevel", slog.LevelError, func(l fastlog.Logger) { l.Debug("msg") }, false},
		{"Info_at_ErrorLevel", slog.LevelError, func(l fastlog.Logger) { l.Info("msg") }, false},
		{"Warn_at_ErrorLevel", slog.LevelError, func(l fastlog.Logger) { l.Warn("msg") }, false},
		{"Error_at_ErrorLevel", slog.LevelError, func(l fastlog.Logger) { l.Error("msg") }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := newCapture()
			logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).LogLevel(tt.logLevel).JSON()

			tt.call(logger)

			hasOutput := len(cw.Lines()) > 0
			if hasOutput != tt.expect {
				t.Errorf("expected output=%v, got output=%v (lines=%d)", tt.expect, hasOutput, len(cw.Lines()))
			}
		})
	}
}

func TestDebugMode(t *testing.T) {
	cw := newCapture()
	logger := fastlog.New().Writer(cw).Caller(false).Timestamp(false).Colors(false).DebugMode(true).JSON()

	logger.Debug("debug msg")

	if len(cw.Lines()) != 1 {
		t.Errorf("expected Debug to produce output with DebugMode(true), got %d lines", len(cw.Lines()))
	}
}

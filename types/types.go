package types

type LogFormat string

const (
	LogFormatJSON    LogFormat = "json"
	LogFormatConsole LogFormat = "console"
	LogFormatLogfmt  LogFormat = "logfmt"
)

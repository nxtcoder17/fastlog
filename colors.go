package fastlog

const ColorReset = "\033[0m"

// Foreground Colors (fg)
const (
	FgBlack   = "\033[38;5;0m"
	FgRed     = "\033[38;5;1m"
	FgGreen   = "\033[38;5;2m"
	FgYellow  = "\033[38;5;3m"
	FgBlue    = "\033[38;5;4m"
	FgMagenta = "\033[38;5;5m"
	FgCyan    = "\033[38;5;6m"
	FgWhite   = "\033[38;5;7m"

	FgBrightBlack   = "\033[38;5;8m"
	FgBrightRed     = "\033[38;5;9m"
	FgBrightGreen   = "\033[38;5;10m"
	FgBrightYellow  = "\033[38;5;11m"
	FgBrightBlue    = "\033[38;5;12m"
	FgBrightMagenta = "\033[38;5;13m"
	FgBrightCyan    = "\033[38;5;14m"
	FgBrightWhite   = "\033[38;5;15m"
)

// Background Colors (bg)
const (
	BgBlack   = "\033[48;5;0m"
	BgRed     = "\033[48;5;1m"
	BgGreen   = "\033[48;5;2m"
	BgYellow  = "\033[48;5;3m"
	BgBlue    = "\033[48;5;4m"
	BgMagenta = "\033[48;5;5m"
	BgCyan    = "\033[48;5;6m"
	BgWhite   = "\033[48;5;7m"

	BgBrightBlack   = "\033[48;5;8m"
	BgBrightRed     = "\033[48;5;9m"
	BgBrightGreen   = "\033[48;5;10m"
	BgBrightYellow  = "\033[48;5;11m"
	BgBrightBlue    = "\033[48;5;12m"
	BgBrightMagenta = "\033[48;5;13m"
	BgBrightCyan    = "\033[48;5;14m"
	BgBrightWhite   = "\033[48;5;15m"
)

const (
	BgDimGray  = "\033[48;5;236m" // darker gray background
	BgDimGreen = "\033[48;5;64m"  // darker gray background
	FgLight    = "\033[38;5;250m" // light foreground

)

const (
	ColorMessage   = FgCyan
	ColorKey       = FgMagenta
	ColorSeparator = FgBlack
)

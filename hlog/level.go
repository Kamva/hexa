package hlog

import "fmt"

// Level can use by all drivers to map to the real level of
// their logger.
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	MessageLevel
	WarnLevel
	ErrorLevel
)

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case MessageLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) CanLog(target Level) bool {
	return l <= target
}

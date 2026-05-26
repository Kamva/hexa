package hlog

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Level can use by all drivers to map to the real level of
// their logger.
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
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
	case ErrorLevel:
		return "error"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) CanLog(target Level) bool {
	return l <= target
}

// LevelFromString converts a string to a log level, returning an error for an
// unknown value instead of panicking.
func LevelFromString(l string) (Level, error) {
	switch l {
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	default:
		return 0, fmt.Errorf("invalid log level: %q", l)
	}
}

func ZapLevel(l Level) zapcore.Level {
	var zl zapcore.Level
	switch l {
	case DebugLevel:
		zl = zapcore.DebugLevel
	case InfoLevel:
		zl = zapcore.InfoLevel
	case WarnLevel:
		zl = zapcore.WarnLevel
	case ErrorLevel:
		zl = zapcore.ErrorLevel
	default:
		panic(fmt.Sprintf("invalid hexa log level: %s", l))
	}
	return zl
}

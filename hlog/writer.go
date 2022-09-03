package hlog

import "io"

// Writer is an io.Writer implementation using Logger.
// Some packages get standard Go log.Logger as their logger. The
// log.Logger gets as writer as param to use for writing logs. we can
// use this writer as the log.Logger's writer.
type Writer struct {
	fn  func(msg string, args ...Field)
	lvl Level
}

// NewWriter creates a new io.Writer using logger.
// Logger param is optional and its default value is
// the global logger.
func NewWriter(l Logger, lvl Level) io.Writer {
	if l == nil {
		l = GlobalLogger()
	}

	var fn func(msg string, fields ...Field)
	switch lvl {
	case DebugLevel:
		fn = l.Debug
	case InfoLevel:
		fn = l.Info
	case WarnLevel:
		fn = l.Warn
	case ErrorLevel:
		fn = l.Error
	default:
		l.Error("can not detect log level to use in Log Writer implementation, choose info", Int("level", int(lvl)))
		fn = l.Info
	}

	return &Writer{
		fn:  fn,
		lvl: lvl,
	}
}

func (w *Writer) Write(in []byte) (int, error) {
	w.fn(string(in))
	return len(in), nil
}

package hlog

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	entry *logrus.Entry
}

func (l *logrusLogger) Core() interface{} {
	return l.entry
}

func (l *logrusLogger) With(ctx hexa.Context, keyValues ...interface{}) hexa.Logger {
	return l.WithFields(keyValues...)
}

func (l *logrusLogger) WithFields(keyValues ...interface{}) hexa.Logger {
	// if key values is not odd, add another item to make it odd.
	if len(keyValues)%2 != 0 {
		keyValues = append(keyValues, errMissingValue)
	}

	fields, _ := gutil.KeyValuesToMap(keyValues...)

	return NewLogrusDriver(l.entry.WithFields(fields))
}

func (l *logrusLogger) Debug(i ...interface{}) {
	l.entry.Debug(i...)
}

func (l *logrusLogger) Info(i ...interface{}) {
	l.entry.Info(i...)
}

func (l *logrusLogger) Message(i ...interface{}) {
	l.entry.Info(i...)
}

func (l *logrusLogger) Warn(i ...interface{}) {
	l.entry.Warn(i...)
}

func (l *logrusLogger) Error(i ...interface{}) {
	l.entry.Error(i...)
}

// NewLogrusDriver return new instance of logrus that implements hexa logger.
// Deprecated
func NewLogrusDriver(logger *logrus.Entry) hexa.Logger {
	return &logrusLogger{
		entry: logger.WithFields(nil),
	}
}

// Assert logrusLogger implements hexa Logger.
var _ hexa.Logger = &logrusLogger{}

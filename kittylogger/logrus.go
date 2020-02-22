package kittylogger

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	entry *logrus.Entry
}

func (l *logrusLogger) WithFields(keyValues ...interface{}) kitty.Logger {
	// if key values is not odd, add another item to make it odd.
	if len(keyValues)%2 != 0 {
		keyValues = append(keyValues, errMissingValue)
	}

	fields, _ := gutil.KeyValuesToMap(keyValues...)

	return NewLogrusDriver(l.entry.WithFields(fields))
}

func (l *logrusLogger) Core() interface{} {
	return l.entry
}

func (l *logrusLogger) Debug(i ...interface{}) {
	l.entry.Debug(i...)
}

func (l *logrusLogger) Info(i ...interface{}) {
	l.entry.Info(i...)
}

func (l *logrusLogger) Warn(i ...interface{}) {
	l.entry.Warn(i...)
}

func (l *logrusLogger) Error(i ...interface{}) {
	l.entry.Error(i...)
}

func (l *logrusLogger) Fatal(i ...interface{}) {
	l.entry.Fatalln(i...)
}

func (l *logrusLogger) Panic(i ...interface{}) {
	l.entry.Panic(i...)
}

// NewLogrusDriver return new instance of logrus that implements kitty entry.
func NewLogrusDriver(logger *logrus.Entry) kitty.Logger {
	return &logrusLogger{entry: logger.WithFields(nil)}
}

// Assert logrusLogger implements kitty Logger.
var _ kitty.Logger = &logrusLogger{}

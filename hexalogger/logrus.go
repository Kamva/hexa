package hexalogger

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

func (l *logrusLogger) With(keyValues ...interface{}) hexa.Logger {
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

func (l *logrusLogger) Error(i ...interface{}) {
	l.entry.Error(i...)
}

func (l *logrusLogger) Fatal(i ...interface{}) {
	l.entry.Fatalln(i...)
}

func (l *logrusLogger) Panic(i ...interface{}) {
	l.entry.Panic(i...)
}

// NewLogrusDriver return new instance of logrus that implements hexa logger.
func NewLogrusDriver(logger *logrus.Entry) hexa.Logger {
	return &logrusLogger{
		entry: logger.WithFields(nil),
	}
}

// Assert logrusLogger implements hexa Logger.
var _ hexa.Logger = &logrusLogger{}

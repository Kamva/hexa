package hlog

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
)

type printerLogger struct {
	level Level
	with  map[string]interface{}
}

func (l *printerLogger) Core() interface{} {
	return fmt.Println
}

func (l *printerLogger) newWith() map[string]interface{} {
	newWith := make(map[string]interface{})
	for k, v := range l.with {
		newWith[k] = v
	}
	return newWith
}

func (l *printerLogger) With(ctx hexa.Context, keyValues ...interface{}) hexa.Logger {
	// if key values is not odd, add another item to make it odd.
	if len(keyValues)%2 != 0 {
		keyValues = append(keyValues, errMissingValue)
	}

	newWith := l.newWith()
	m, _ := gutil.KeyValuesToMap(keyValues...)
	gutil.ExtendMap(newWith, m, true)

	newLogger := NewPrinterDriver(PrinterOptions{}).(*printerLogger)
	newLogger.with = newWith
	return newLogger
}

func (l *printerLogger) WithFields(args ...interface{}) hexa.Logger {
	return l.With(nil, args...)
}

func (l *printerLogger) log(level Level, i ...interface{}) {
	if l.level.CanLog(level) {
		fmt.Println(fmt.Sprintf("%s: ", level), l.with, i)
	}
}

func (l *printerLogger) Debug(i ...interface{}) {
	l.log(DebugLevel, i...)
}

func (l *printerLogger) Info(i ...interface{}) {
	l.log(InfoLevel, i...)
}

func (l *printerLogger) Message(i ...interface{}) {
	l.log(MessageLevel, i...)
}

func (l *printerLogger) Warn(i ...interface{}) {
	l.log(WarnLevel, i...)
}

func (l *printerLogger) Error(i ...interface{}) {
	l.log(ErrorLevel, i...)
}

type PrinterOptions struct {
}

// NewPrinterDriver returns new instance of hexa logger
// with printer driver.
// Note: printer logger driver is just for test purpose.
// dont use it in production.
func NewPrinterDriver(o PrinterOptions) hexa.Logger {
	return NewPrinterDriverWith(DebugLevel)
}

func NewPrinterDriverWith(l Level) hexa.Logger {
	return &printerLogger{level: l, with: map[string]interface{}{}}
}

// Assert printerLogger implements hexa Logger.
var _ hexa.Logger = &printerLogger{}

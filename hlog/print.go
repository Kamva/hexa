package hlog

import (
	"fmt"
	"github.com/kamva/hexa"
)

type printerLogger struct {
	level Level
	with  []Field
}

func (l *printerLogger) Core() interface{} {
	return fmt.Println
}

func (l *printerLogger) newWith() []Field {
	dst := make([]Field, len(l.with))
	copy(l.with, dst)
	return dst
}

func (l *printerLogger) WithCtx(ctx hexa.Context, args ...Field) hexa.Logger {
	newWith := l.newWith()
	newWith = append(newWith, args...)

	newLogger := NewPrinterDriver().(*printerLogger)
	newLogger.with = newWith
	return newLogger
}

func (l *printerLogger) With(args ...Field) hexa.Logger {
	return l.WithCtx(nil, args...)
}

func (l *printerLogger) WithFunc(f hexa.LogFunc) hexa.Logger {
	return f(l)
}

func (l *printerLogger) log(level Level, msg string, args ...Field) {
	ll := l.With(args...).(*printerLogger)
	if l.level.CanLog(level) {
		fmt.Println(fmt.Sprintf("%s: ", level), fieldsToMap(ll.with...), msg)
	}
}

func (l *printerLogger) Debug(msg string, args ...Field) {
	l.log(DebugLevel, msg, args...)
}

func (l *printerLogger) Info(msg string, args ...Field) {
	l.log(InfoLevel, msg, args...)
}

func (l *printerLogger) Message(msg string, args ...Field) {
	l.log(MessageLevel, msg, args...)
}

func (l *printerLogger) Warn(msg string, args ...Field) {
	l.log(WarnLevel, msg, args...)
}

func (l *printerLogger) Error(msg string, args ...Field) {
	l.log(ErrorLevel, msg, args...)
}

// NewPrinterDriver returns new instance of hexa logger
// with printer driver.
// Note: printer logger driver is just for test purpose.
// dont use it in production.
func NewPrinterDriver() hexa.Logger {
	return NewPrinterDriverWith(DebugLevel)
}

func NewPrinterDriverWith(l Level) hexa.Logger {
	return &printerLogger{level: l, with: make([]Field, 0)}
}

// Assert printerLogger implements hexa Logger.
var _ hexa.Logger = &printerLogger{}

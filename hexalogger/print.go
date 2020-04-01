package hexalogger

import (
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

type printerLogger struct {
	with map[string]interface{}
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

	newLogger := NewPrinterDriver().(*printerLogger)
	newLogger.with = newWith
	return newLogger
}

func (l *printerLogger) WithFields(args ...interface{}) hexa.Logger {
	return l.With(nil, args...)
}

func (l *printerLogger) Debug(i ...interface{}) {
	fmt.Println("Debug: ", l.with, i)
}

func (l *printerLogger) Info(i ...interface{}) {
	fmt.Println("Info: ", l.with, i)
}

func (l *printerLogger) Message(i ...interface{}) {
	fmt.Println("Message: ", l.with, i)
}

func (l *printerLogger) Warn(i ...interface{}) {
	fmt.Println("Warn: ", l.with, i)
}

func (l *printerLogger) Error(i ...interface{}) {
	fmt.Println("Error: ", l.with, i)
}

// NewPrinterDriver return new instance of hexa logger
// with printer driver.
// Note: printer logger driver is just for test purpose.
// dont use it in production.
func NewPrinterDriver() hexa.Logger {
	return &printerLogger{with: map[string]interface{}{}}
}

// Assert printerLogger implements hexa Logger.
var _ hexa.Logger = &printerLogger{}

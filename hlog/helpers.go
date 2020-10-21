package hlog

import (
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// WithTrace add error trace (if exists) to the log
func WithTrace(l hexa.Logger, err error) hexa.Logger {
	if stack := tracer.StackAsString(err); stack != "" {
		return l.WithFields(hexa.TracedStackLogKey, fmt.Sprintf("%+v", stErr))
	}
	return l
}

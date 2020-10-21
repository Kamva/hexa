package hlog

import (
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// WithErrStack add error stack (if exists) to the log
func WithErrStack(l hexa.Logger, err error) hexa.Logger {
	if stack := tracer.StackAsString(err); stack != "" {
		return l.WithFields(hexa.ErrorStackLogKey, fmt.Sprintf("%+v", stack))
	}
	return l
}

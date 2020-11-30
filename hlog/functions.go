package hlog

import (
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// ErrStack add error stack (if exists) to the log
func ErrStackF(err error) hexa.LogFunc {
	return func(l hexa.Logger) hexa.Logger {
		if stack := tracer.StackAsString(err); stack != "" {
			return l.WithFields(String(hexa.ErrorStackLogKey, fmt.Sprintf("%+v", stack)))
		}
		return l
	}
}

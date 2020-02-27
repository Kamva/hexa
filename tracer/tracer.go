//  tracer package add stack trace to the errors.
//  - Use unwrap function from standard errors package to unwrap error.
//  - Use Is function from standard errors package to check error is expected error or no.
//
package tracer

import "github.com/pkg/errors"

// stack represents a stack of program counters.
type (
	// traceErr is the error struct that contain trace of error.
	StackTracer interface {
		StackTrace() errors.StackTrace
	}
)

// Trace function check if error contains trace, so
// return it, otherwise add stacktrace to the error.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackTracer); ok {
		return err
	}

	return errors.WithStack(err)
}

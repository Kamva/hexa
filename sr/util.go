package sr

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// ShutdownBySignals shutdown service registry services by listening to the os signals.
func ShutdownBySignals(sr hexa.ServiceRegistry, timeout time.Duration, signals ...os.Signal) error {
	if len(signals) == 0 { // set default signals.
		signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	gutil.WaitForSignals(signals...)

	ctx := context.Background()
	if timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	return tracer.Trace(sr.Shutdown(ctx))
}

// ShutdownWithTimeout shuts down with the provided timeout.
func ShutdownWithTimeout(r hexa.ServiceRegistry, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return tracer.Trace(r.Shutdown(ctx))
}

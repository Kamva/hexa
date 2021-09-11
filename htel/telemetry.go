package htel

import (
	"context"

	"github.com/kamva/tracer"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetry is just a wrapper for openTelemetry services to
// implement hexa services(to shutdown... them).
type OpenTelemetry interface {
	TracerProvider() trace.TracerProvider
}

func NewOpenTelemetry(tp *tracesdk.TracerProvider) OpenTelemetry {
	return &openTelemetry{tp: tp}
}

type openTelemetry struct {
	tp *tracesdk.TracerProvider
}

func (t *openTelemetry) TracerProvider() trace.TracerProvider {
	return t.tp
}

func (t *openTelemetry) Shutdown(c context.Context) error {
	// Shutdown other open telemetry instances here.
	return tracer.Trace(t.tp.Shutdown(c))
}

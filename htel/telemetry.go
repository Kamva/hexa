package htel

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"go.opentelemetry.io/otel/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// tracerProvider is a wrapper around the tracesdk.TracerProvider to be sure implements a hexa
// service to shutdown the TracerProvider.
type tracerProvider struct {
	*tracesdk.TracerProvider
}

func (t *tracerProvider) Shutdown(ctx context.Context) error {
	return tracer.Trace(t.TracerProvider.Shutdown(ctx))
}

func NewTracerProvider(tp *tracesdk.TracerProvider) trace.TracerProvider {
	return &tracerProvider{tp}
}

// OpenTelemetry is just a wrapper for openTelemetry services to
// implement hexa services(to shutdown... them).
type OpenTelemetry interface {
	TracerProvider() trace.TracerProvider
	MeterProvider() metric.MeterProvider
}

type openTelemetry struct {
	tp trace.TracerProvider
	mp metric.MeterProvider
}

func NewOpenTelemetry(tp trace.TracerProvider, mp metric.MeterProvider) OpenTelemetry {
	return &openTelemetry{tp: tp, mp: mp}
}

func (t *openTelemetry) TracerProvider() trace.TracerProvider {
	return t.tp
}

func (t *openTelemetry) MeterProvider() metric.MeterProvider {
	return t.mp
}

var _ hexa.Shutdownable = &tracerProvider{}

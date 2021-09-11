package htel

import (
	"context"
	"testing"

	"github.com/kamva/tracer"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func traceProvider() (*tracesdk.TracerProvider, error) {
	ex, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(ex),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("abc"),
			attribute.String("environment", "test"),
			attribute.Int64("ID", 1),
		)),
	), nil
}

func TestHexaPropagator_ExtractEmptyContext(t *testing.T) {
	otelPropagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	hp := NewHexaPropagator(otelPropagator)

	m, err := hp.Inject(context.Background())
	assert.Nil(t, err)
	assert.Equalf(t, 0, len(m), "extracted map must be empty")
}

func TestHexaPropagator_Extract(t *testing.T) {
	otelPropagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	hp := NewHexaPropagator(otelPropagator)

	p, _ := traceProvider()
	ctx, _ := p.Tracer("test").Start(context.Background(), "abc")

	m, err := hp.Inject(ctx)
	assert.Nil(t, err)
	assert.Equalf(t, 1, len(m), "must insert one item to the result map")
}

// TODO: write more tests.
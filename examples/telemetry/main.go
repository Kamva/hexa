package main

import (
	"context"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "hexa-demo"
	environment = "dev"
	id          = 1
)

var ot htel.OpenTelemetry

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
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	), nil
}

func main() {
	r := sr.New()

	tp, err := traceProvider()
	gutil.PanicErr(err)
	otel.SetTracerProvider(tp) // set global tracer too (we don't need to it although)

	ot = htel.NewOpenTelemetry(tp,metric.NewNoopMeterProvider())
	r.Register("open_telemetry", ot) // register openTelemetry as a service in the hexa service registry
	defer r.Shutdown(context.Background())

	t := ot.TracerProvider().Tracer("main-component")
	ctx, span := t.Start(context.Background(), "foo")
	defer span.End()

	bar(ctx)
}

func bar(ctx context.Context) {
	t := ot.TracerProvider().Tracer("bar-component")
	_, span := t.Start(ctx, "bar")
	defer span.End()

	// Do something here...
}

package htel

import (
	"context"
	"testing"

	"github.com/kamva/hexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// noopProcessor is a span processor so the SDK TracerProvider shuts down
// cleanly (in this SDK version Shutdown errors when none is registered).
type noopProcessor struct{}

func (noopProcessor) OnStart(context.Context, tracesdk.ReadWriteSpan) {}
func (noopProcessor) OnEnd(tracesdk.ReadOnlySpan)                     {}
func (noopProcessor) Shutdown(context.Context) error                  { return nil }
func (noopProcessor) ForceFlush(context.Context) error                { return nil }

func newSDKTracerProvider() *tracesdk.TracerProvider {
	return tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(noopProcessor{}))
}

func TestHexaCarrier(t *testing.T) {
	c := make(HexaCarrier)
	c.Set("k", "v")
	c.Set("k2", "v2")

	assert.Equal(t, "v", c.Get("k"))
	assert.Equal(t, "", c.Get("missing"))
	assert.ElementsMatch(t, []string{"k", "k2"}, c.Keys())
}

func TestHexaPropagator_InjectExtractRoundTrip(t *testing.T) {
	tid, err := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	require.NoError(t, err)
	sid, err := trace.SpanIDFromHex("0102030405060708")
	require.NoError(t, err)

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	p := NewHexaPropagator(propagation.TraceContext{})

	m, err := p.Inject(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, m["traceparent"])

	got, err := p.Extract(context.Background(), m)
	require.NoError(t, err)
	gotSC := trace.SpanContextFromContext(got)
	assert.Equal(t, tid, gotSC.TraceID())
	assert.Equal(t, sid, gotSC.SpanID())
}

func TestTracerProvider_Shutdownable(t *testing.T) {
	tp := NewTracerProvider(newSDKTracerProvider())

	sd, ok := tp.(hexa.Shutdownable)
	require.True(t, ok)
	assert.NoError(t, sd.Shutdown(context.Background()))
}

func TestOpenTelemetry_Getters(t *testing.T) {
	tp := NewTracerProvider(newSDKTracerProvider())
	ot := NewOpenTelemetry(tp, nil)

	assert.Equal(t, tp, ot.TracerProvider())
	assert.Nil(t, ot.MeterProvider())
}

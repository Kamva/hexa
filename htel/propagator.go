package htel

import (
	"context"

	"github.com/kamva/hexa"
	"go.opentelemetry.io/otel/propagation"
)

// HexaCarrier is an open telemetry carrier to get & give exported open telemetry data.
type HexaCarrier map[string][]byte

func (hc HexaCarrier) Get(key string) string {
	return string(hc[key])
}

func (hc HexaCarrier) Set(key string, value string) {
	hc[key] = []byte(value)
}

func (hc HexaCarrier) Keys() []string {
	keys := make([]string, len(hc))
	i := 0
	for k, _ := range hc {
		keys[i] = k
		i++
	}
	return keys
}

// hexaPropagator propagates opentelemetry data as a hexa propagator.
type hexaPropagator struct {
	p propagation.TextMapPropagator
}

func NewHexaPropagator(p propagation.TextMapPropagator) hexa.ContextPropagator {
	return &hexaPropagator{p: p}
}

func (o *hexaPropagator) Inject(ctx context.Context) (map[string][]byte, error) {
	carrier := make(HexaCarrier)
	o.p.Inject(ctx, carrier)
	return carrier, nil
}

func (o *hexaPropagator) Extract(ctx context.Context, m map[string][]byte) (context.Context, error) {
	return o.p.Extract(ctx, HexaCarrier(m)), nil
}

var _ hexa.ContextPropagator = &hexaPropagator{}
var _ propagation.TextMapCarrier = &HexaCarrier{}

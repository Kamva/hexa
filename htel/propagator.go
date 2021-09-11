package htel

import (
	"context"

	"github.com/kamva/hexa"
	"go.opentelemetry.io/otel/propagation"
)

// hexaPropagator propagates opentelemetry data as a hexa propagator.
type hexaPropagator struct {
	p propagation.TextMapPropagator
}

func NewHexaPropagator(p propagation.TextMapPropagator) hexa.ContextPropagator {
	return &hexaPropagator{p: p}
}

func (o *hexaPropagator) Inject(ctx context.Context) (map[string][]byte, error) {
	carrier := make(HexaCarrier)
	o.p.Inject(ctx, HexaCarrier(carrier))
	return carrier, nil
}

func (o *hexaPropagator) Extract(ctx context.Context, m map[string][]byte) (context.Context, error) {
	return o.p.Extract(ctx, HexaCarrier(m)), nil
}

var _ propagation.TextMapCarrier = &HexaCarrier{}
var _ hexa.ContextPropagator = &hexaPropagator{}

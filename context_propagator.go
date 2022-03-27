package hexa

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/tracer"
)

type ContextPropagator interface {
	// Inject injects values from context into a map and returns the map.
	Inject(context.Context) (map[string][]byte, error)

	// Extract extract values from map and insert to the context.
	Extract(context.Context, map[string][]byte) (context.Context, error)
}

var propagatingContextKeys = []contextKey{
	ctxKeyCorrelationId,
	ctxKeyLocale,
	ctxKeyUser,
}

// keysPropagator get a list of keys and propagate that key,values from context.
// all values in the context for these keys must be string.
type keysPropagator struct {
	keys   []string
	strict bool
}

// defaultContextPropagator propagate the default implementation of the Hexa context.
// You can use it as one of hexa propagators to propagate hexa context itself across
// microservices.
type defaultContextPropagator struct {
	up         UserPropagator
	logger     Logger
	translator Translator
}

// multiPropagator get multiple propagator and return as one
// propagator which you can use as your propagator.
type multiPropagator struct {
	propagators []ContextPropagator
}

func (p *multiPropagator) Inject(c context.Context) (map[string][]byte, error) {
	finalMap := make(map[string][]byte)
	for _, p := range p.propagators {
		m, err := p.Inject(c)
		if err != nil {
			return nil, tracer.Trace(err)
		}
		extendBytesMap(finalMap, m, true)
	}
	return finalMap, nil
}

func (p *multiPropagator) Extract(c context.Context, m map[string][]byte) (context.Context, error) {
	var err error
	for _, p := range p.propagators {
		c, err = p.Extract(c, m)
		if err != nil {
			return nil, tracer.Trace(err)
		}
	}
	return c, nil
}

func (p *multiPropagator) AddPropagator(propagator ContextPropagator) {
	p.propagators = append(p.propagators, propagator)
}

func (p *defaultContextPropagator) Inject(c context.Context) (map[string][]byte, error) {
	// just get local, correlation_id  and user
	m := make(map[string][]byte)
	m[string(ctxKeyCorrelationId)] = []byte(CtxCorrelationId(c))
	m[string(ctxKeyLocale)] = []byte(CtxLocale(c))

	// user
	user := CtxUser(c)
	if user != nil {
		uBytes, err := p.up.ToBytes(user)
		if err != nil {
			return nil, tracer.Trace(err)
		}
		m[string(ctxKeyUser)] = uBytes
	}

	return m, nil
}

func (p *defaultContextPropagator) Extract(c context.Context, m map[string][]byte) (context.Context, error) {
	for _, k := range propagatingContextKeys {
		if _, ok := m[string(k)]; !ok {
			return nil, tracer.Trace(fmt.Errorf("key %s not found in map", k))
		}
	}

	user, err := p.up.FromBytes(m[string(ctxKeyUser)])
	if err != nil { // Another option would be ignoring the nil user and continue...
		return nil, tracer.Trace(err)
	}

	return NewContext(c, ContextParams{
		Request:        nil,
		CorrelationId:  string(m[string(ctxKeyCorrelationId)]),
		Locale:         string(m[string(ctxKeyLocale)]),
		User:           user,
		BaseLogger:     p.logger,
		BaseTranslator: p.translator,
		Store:          newStore(),
	}), nil
}

func (p *keysPropagator) Inject(c context.Context) (map[string][]byte, error) {
	m := make(map[string][]byte)

	for _, k := range p.keys {
		val, ok := c.Value(k).(string)
		if !ok {
			return nil, tracer.Trace(fmt.Errorf("type of value for %s key is not string", k))
		}
		m[k] = []byte(val)
	}

	return m, nil
}

func (p *keysPropagator) Extract(c context.Context, m map[string][]byte) (context.Context, error) {
	for _, k := range p.keys {
		v, ok := m[k]
		if !ok {
			if p.strict {
				return nil, tracer.Trace(fmt.Errorf("value for key %s does not exist", k))
			}
			continue
		}
		c = context.WithValue(c, k, string(v))
	}
	return c, nil
}

func NewMultiPropagator(propagators ...ContextPropagator) ContextPropagator {
	return &multiPropagator{propagators: propagators}
}

// NewContextPropagator returns new context propagator to propagate
// the Hexa context itself.
func NewContextPropagator(l Logger, t Translator) ContextPropagator {
	return &defaultContextPropagator{up: NewUserPropagator(), logger: l, translator: t}
}

func NewKeysPropagator(keys []string, strict bool) ContextPropagator {
	return &keysPropagator{keys: keys, strict: strict}
}

// WithPropagator add another propagator to ourself implemented multiPropagator.
func WithPropagator(multi ContextPropagator, p ContextPropagator) error {
	multiP, ok := multi.(*multiPropagator)
	if !ok {
		msg := "propagator is not multi propagator, we can not add another propagator to it."
		return tracer.Trace(errors.New(msg))
	}
	multiP.AddPropagator(p)
	return nil
}

var _ ContextPropagator = &multiPropagator{}
var _ ContextPropagator = &defaultContextPropagator{}
var _ ContextPropagator = &keysPropagator{}

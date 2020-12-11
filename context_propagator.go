package hexa

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/tracer"
)

type ContextPropagator interface {
	// Extract extracts values from context and add to the map
	Extract(context.Context) (map[string][]byte, error)

	// Inject injects values from map to the context.
	Inject(map[string][]byte, context.Context) (context.Context, error)
}

var propagatingContextKeys = []string{
	ContextKeyCorrelationID,
	ContextKeyLocale,
	ContextKeyUser,
}

// keysPropagator get propagators of keys which should propagate from context.
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

type multiPropagator struct {
	propagators []ContextPropagator
}

func (p *multiPropagator) Extract(c context.Context) (map[string][]byte, error) {
	finalMap := make(map[string][]byte)
	for _, p := range p.propagators {
		m, err := p.Extract(c)
		if err != nil {
			return nil, tracer.Trace(err)
		}
		extendBytesMap(finalMap, m, true)
	}
	return finalMap, nil
}

func (p *multiPropagator) Inject(m map[string][]byte, c context.Context) (context.Context, error) {
	var err error
	for _, p := range p.propagators {
		c, err = p.Inject(m, c)
		if err != nil {
			return nil, tracer.Trace(err)
		}
	}
	return c, nil
}

func (p *multiPropagator) AddPropagator(propagator ContextPropagator) {
	p.propagators = append(p.propagators, propagator)
}

func (p *defaultContextPropagator) Extract(c context.Context) (map[string][]byte, error) {
	// just extract local, correlation_id  and user
	m := make(map[string][]byte)
	m[ContextKeyCorrelationID] = []byte(c.Value(ContextKeyCorrelationID).(string))
	m[ContextKeyLocale] = []byte(c.Value(ContextKeyLocale).(string))

	// user
	user := c.Value(ContextKeyUser).(User)
	uBytes, err := p.up.ToBytes(user)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	m[ContextKeyUser] = uBytes

	return m, nil
}

func (p *defaultContextPropagator) Inject(m map[string][]byte, c context.Context) (context.Context, error) {
	for _, k := range propagatingContextKeys {
		if _, ok := m[k]; !ok {
			return nil, tracer.Trace(fmt.Errorf("key %s not found in map", k))
		}
	}
	user, err := p.up.FromBytes(m[ContextKeyUser])
	if err != nil {
		return nil, tracer.Trace(err)
	}

	// Set context values:
	c = contextWithParams(c, ContextParams{
		Request:       nil,
		CorrelationId: string(m[ContextKeyCorrelationID]),
		Locale:        string(m[ContextKeyLocale]),
		User:          user,
		Logger:        p.logger,
		Translator:    p.translator,
	})
	return c, nil
}

func (p *keysPropagator) Extract(c context.Context) (map[string][]byte, error) {
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

func (p *keysPropagator) Inject(m map[string][]byte, c context.Context) (context.Context, error) {
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

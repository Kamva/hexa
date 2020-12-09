package hexa

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kamva/gutil"
	"github.com/stretchr/testify/assert"
)

func assertImportedContextWithParams(t *testing.T, ctx Context, params *ContextParams) {
	assert.Nil(t,ctx.Request())
	assert.Equal(t, params.CorrelationId, ctx.CorrelationID())
	assert.Equal(t, params.Locale, ctx.Value(ContextKeyLocale))
	assert.Equal(t, params.User, ctx.User())
	assert.NotNil(t, ctx.Logger())
	assert.NotNil(t, ctx.Translator())
}


func TestDefaultContextPropagator_Extract(t *testing.T) {
	context, params := newTestContext()
	l := &emptyLogger{}
	translator := &emptyTranslator{}
	p := NewContextPropagator(l, translator)

	uBytes, err := json.Marshal(params.User.MetaData())
	gutil.PanicErr(err)
	result := map[string][]byte{
		ContextKeyCorrelationID: []byte(params.CorrelationId),
		ContextKeyLocale:        []byte(params.Locale),
		ContextKeyUser:          uBytes,
	}
	m, err := p.Extract(context)
	assert.Nil(t, err)
	assert.Equal(t, result, m)
}

func TestDefaultContextPropagator_Inject(t *testing.T) {
	_, params := newTestContext()
	l := &emptyLogger{}
	translator := &emptyTranslator{}
	p := NewContextPropagator(l, translator)

	uBytes, err := json.Marshal(params.User.MetaData())
	gutil.PanicErr(err)
	payload := map[string][]byte{
		ContextKeyCorrelationID: []byte(params.CorrelationId),
		ContextKeyLocale:        []byte(params.Locale),
		ContextKeyUser:          uBytes,
	}
	result, err := p.Inject(payload, context.Background())
	assert.Nil(t, err)
	ctx, err := NewContextFromRawContext(result)
	if !assert.Nil(t, err) {
		return
	}

	assertImportedContextWithParams(t, ctx, params)
}

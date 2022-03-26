package hexa

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kamva/gutil"
	"github.com/stretchr/testify/assert"
)

func assertImportedContextWithParams(t *testing.T, ctx context.Context, params *ContextParams) {
	assert.Nil(t, CtxRequest(ctx))
	assert.Equal(t, params.CorrelationId, CtxCorrelationId(ctx))
	assert.Equal(t, params.Locale, CtxLocale(ctx))
	assert.Equal(t, params.User, CtxUser(ctx))
	assert.NotNil(t, CtxLogger(ctx))
	assert.NotNil(t, CtxTranslator(ctx))
}

func TestDefaultContextPropagator_Extract(t *testing.T) {
	context, params := newTestContext()
	l := &emptyLogger{}
	translator := &emptyTranslator{}
	p := NewContextPropagator(l, translator)

	uBytes, err := json.Marshal(params.User.MetaData())
	gutil.PanicErr(err)
	result := map[string][]byte{
		string(ctxKeyCorrelationId): []byte(params.CorrelationId),
		string(ctxKeyLocale):        []byte(params.Locale),
		string(ctxKeyUser):          uBytes,
	}
	m, err := p.Inject(context)
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
		string(ctxKeyCorrelationId): []byte(params.CorrelationId),
		string(ctxKeyLocale):        []byte(params.Locale),
		string(ctxKeyUser):          uBytes,
	}
	ctx, err := p.Extract(context.Background(), payload)
	assert.Nil(t, err)
	assertImportedContextWithParams(t, ctx, params)
}

package hexa

import (
	"context"
	"net/http"
	"testing"

	"github.com/kamva/gutil"
	"github.com/stretchr/testify/assert"
)

func newTestContext() (context.Context, *ContextParams) {
	r := gutil.Must(http.NewRequest("POST", "http://a.com", nil)).(*http.Request)
	params := ContextParams{
		Request:        r,
		CorrelationId:  "abc",
		Locale:         "def",
		User:           NewGuest(),
		BaseLogger:     &emptyLogger{},
		BaseTranslator: &emptyTranslator{},
	}
	return NewContext(nil, params), &params
}

func assertContextWithParams(t *testing.T, ctx context.Context, params *ContextParams) {
	assert.Equal(t, params.Request, CtxRequest(ctx))
	assert.Equal(t, params.CorrelationId, CtxCorrelationId(ctx))
	assert.Equal(t, params.Locale, CtxLocale(ctx))
	assert.Equal(t, params.User, CtxUser(ctx))
	assert.NotNil(t, CtxLogger(ctx))
	assert.NotNil(t, CtxTranslator(ctx))
}

func TestNewContext(t *testing.T) {
	ctx, params := newTestContext()
	if !assert.NotNil(t, ctx) {
		return
	}
	assertContextWithParams(t, ctx, params)
}

func TestWithUser(t *testing.T) {
	ctx, params := newTestContext()
	assert.Equal(t, CtxUser(ctx), params.User)
	newUser := NewGuest()
	newCtx := WithUser(ctx, newUser)
	assert.Equal(t, CtxUser(ctx), params.User)

	assert.Equal(t, CtxUser(ctx), newUser)

	// assert old context is good:
	assertContextWithParams(t, ctx, params)

	params.User = newUser
	assertContextWithParams(t, newCtx, params)
}

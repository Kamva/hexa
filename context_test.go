package hexa

import (
	"context"
	"net/http"
	"testing"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/hlog"
	"github.com/stretchr/testify/assert"
)

func newTestContext() (context.Context, *ContextParams) {
	r := gutil.Must(http.NewRequest("POST", "http://a.com", nil)).(*http.Request)
	params := ContextParams{
		Request:        r,
		CorrelationId:  "abc",
		Locale:         "def",
		User:           NewGuest(),
		BaseLogger:     hlog.GlobalLogger(),
		BaseTranslator: &emptyTranslator{},
	}
	return NewContext(context.Background(), params), &params
}

func assertContextWithParams(ctx context.Context, t *testing.T, params *ContextParams) {
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
	assertContextWithParams(ctx, t, params)
}

func TestWithUser(t *testing.T) {
	ctx, params := newTestContext()
	assert.Equal(t, CtxUser(ctx), params.User)
	newUser := NewGuest()
	newCtx := WithUser(ctx, newUser)
	assert.Equal(t, CtxUser(ctx), params.User)

	assert.Equal(t, CtxUser(ctx), newUser)

	// assert old context is good:
	assertContextWithParams(ctx, t, params)

	params.User = newUser
	assertContextWithParams(newCtx, t, params)
}

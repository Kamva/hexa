package hexa

import (
	"net/http"
	"testing"

	"github.com/kamva/gutil"
	"github.com/stretchr/testify/assert"
)

func newTestContext() (Context, *ContextParams) {
	r := gutil.Must(http.NewRequest("POST", "http://a.com", nil)).(*http.Request)
	params := ContextParams{
		Request:       r,
		CorrelationId: "abc",
		Locale:        "def",
		User:          NewGuest(),
		Logger:        &emptyLogger{},
		Translator:    &emptyTranslator{},
	}
	return NewContext(nil, params), &params
}

func assertContextWithParams(t *testing.T, ctx Context, params *ContextParams) {
	assert.Equal(t, params.Request, ctx.Request())
	assert.Equal(t, params.CorrelationId, ctx.CorrelationID())
	assert.Equal(t, params.Locale, ctx.Value(ContextKeyLocale))
	assert.Equal(t, params.User, ctx.User())
	assert.NotNil(t, ctx.Logger())
	assert.NotNil(t, ctx.Translator())
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
	assert.Equal(t, ctx.User(), params.User)
	newUser := NewGuest()
	newCtx := WithUser(ctx, newUser)
	assert.Equal(t, ctx.User(), params.User)

	assert.Equal(t, newCtx.User(), newUser)

	// assert old context is good:
	assertContextWithParams(t, ctx, params)

	params.User = newUser
	assertContextWithParams(t, newCtx, params)

}

package hexa

import (
	"net/http"
	"testing"

	"github.com/kamva/gutil"
	"github.com/stretchr/testify/assert"
)

func getNewContext() (Context, *ContextParams) {
	r := gutil.Must(http.NewRequest("POST", "http://a.com", nil)).(*http.Request)
	params := ContextParams{
		Request:       r,
		CorrelationId: "abc",
		Locale:        "def",
		User:          NewGuest(),
		Logger:        &emptyLogger{},
		Translator:    &emptyTranslator{},
	}
	return NewContext(params), &params
}

func TestNewContext(t *testing.T) {
	ctx, params := getNewContext()
	if !assert.NotNil(t, ctx) {
		return
	}

	assert.Equal(t, params.Request, ctx.Request())
	assert.Equal(t, params.CorrelationId, ctx.CorrelationID())
	assert.Equal(t, params.Locale, ctx.Value(ContextKeyLocale))
	assert.Equal(t, params.User, ctx.User())
	assert.NotNil(t, ctx.Logger())
	assert.NotNil(t, ctx.Translator())
}

func TestWithUser(t *testing.T) {
	ctx, params := getNewContext()
	assert.Equal(t, ctx.User(), params.User)
	newUser := NewGuest()
	newCtx := WithUser(ctx, newUser)
	assert.Equal(t, ctx.User(), params.User)

	assert.Equal(t, newCtx.User(), newUser)
	assert.Equal(t, params.Request, ctx.Request())
	assert.Equal(t, params.CorrelationId, ctx.CorrelationID())
	assert.Equal(t, params.Locale, ctx.Value(ContextKeyLocale))
	assert.Equal(t, params.User, ctx.User())
	assert.NotNil(t, ctx.Logger())
	assert.NotNil(t, ctx.Translator())
}

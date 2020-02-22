package kitty

import (
	"context"
)

// Context is the kitty context to use in services.
type Context interface {
	context.Context

	// RequestID returns the request id.
	RequestID() string
	// User returns the user
	User() User

	// Logger returns the kitty logger customized for specific request.
	Logger() Logger

	// Translator returns the translator localized relative to the users request.
	Translator() Translator
}

type defaultContext struct {
	context.Context
	requestID  string
	user       User
	logger     Logger
	translator Translator
}

func (c defaultContext) RequestID() string {
	return c.requestID
}

func (c defaultContext) User() User {
	return c.user
}

func (c defaultContext) Logger() Logger {
	return c.logger
}

func (c defaultContext) Translator() Translator {
	return c.translator
}

// NewCtx returns new kitty Context.
func NewCtx(requestID string, user User, logger Logger, translator Translator) Context {
	return &defaultContext{
		Context:    context.Background(),
		requestID:  requestID,
		user:       user,
		logger:     logger,
		translator: translator,
	}
}

// Assert defaultContext implements the kitty Context.
var _ Context = &defaultContext{}

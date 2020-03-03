package kitty

import (
	"context"
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/tracer"
)

type UserFinder = func(id interface{}) (User, error)

// Context is the kitty context to use in services.
type Context interface {
	context.Context

	// RequestID returns the request id.
	RequestID() string

	// Correlation returns the request correlation id.
	CorrelationID() string

	// User returns the user
	User() User

	// Logger returns the kitty logger customized for specific request.
	Logger() Logger

	// Translator returns the translator localized relative to the users request.
	Translator() Translator

	// ToMap method convert context to map to export and import.
	ToMap() Map
}

type defaultContext struct {
	context.Context
	locale        string
	requestID     string
	correlationID string
	user          User
	logger        Logger
	translator    Translator
}

// exportedCtx use export and import Context.
type exportedCtx struct {
	Locale        string      `json:"locale"`
	RequestID     string      `json:"request_id"`
	CorrelationID string      `json:"correlation_id"`
	UserId        interface{} `json:"user_id"`
}

func (e exportedCtx) validate() error {
	if e.RequestID == "" || e.CorrelationID == "" || e.UserId == "" {
		return tracer.Trace(errors.New("exported data is invalid"))
	}

	return nil
}

func (c defaultContext) RequestID() string {
	return c.requestID
}

func (c defaultContext) CorrelationID() string {
	return c.correlationID
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

func (c defaultContext) ToMap() Map {
	return gutil.StructToMap(exportedCtx{
		Locale:        c.locale,
		RequestID:     c.RequestID(),
		CorrelationID: c.CorrelationID(),
		UserId:        c.User().Identifier().Val(),
	})
}

// NewCtx returns new kitty Context.
// locale syntax is just same as HTTP Accept-Language header.
func NewCtx(requestID, correlationID string, locale string, user User, logger Logger, translator Translator) Context {
	logger = tuneCtxLogger(requestID, correlationID, user, logger)
	translator = tuneCtxTranslator(locale, translator)

	return &defaultContext{
		Context:       context.Background(),
		locale:        locale,
		requestID:     requestID,
		correlationID: correlationID,
		user:          user,
		logger:        logger,
		translator:    translator,
	}
}

// CtxFromMap generate new context form the exported map by
// another context (by using ToMap function on the context).
func CtxFromMap(m Map, uf UserFinder, l Logger, t Translator) (Context, error) {
	e := exportedCtx{}
	err := gutil.MapToStruct(m, &e)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	if err := e.validate(); err != nil {
		return nil, tracer.Trace(err)
	}

	u := NewGuestUser()

	if e.UserId != guestUserID {
		u, err = uf(e.UserId)
		if err != nil {
			return nil, tracer.Trace(err)
		}
	}

	return NewCtx(e.RequestID, e.CorrelationID, e.Locale, u, l, t), nil
}

// tuneLogger function tune the logger for each context.
func tuneCtxLogger(requestID string, correlationID string, u User, logger Logger) Logger {

	logger = logger.WithFields(
		"__guest__", u.IsGuest(),
		"__user_id__", u.Identifier().String(),
		"__username__", u.GetUsername(),
		"__request_id__", requestID,
		"__correlation_id__", correlationID,
	)

	return logger
}

// tuneCtxTranslator localize translator for each context.
func tuneCtxTranslator(locale string, t Translator) Translator {
	if locale != "" {
		return t.Localize(locale)
	}

	return t.Localize()
}

// Assert defaultContext implements the kitty Context.
var _ Context = &defaultContext{}

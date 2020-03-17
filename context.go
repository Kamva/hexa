package hexa

import (
	"context"
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/tracer"
	"net/http"
)

type UserFinder = func(id interface{}) (User, error)

// Context is the hexa context to use in services.
type Context interface {
	context.Context

	// Request returns the current request and can be nil for not http requests.
	Request() *http.Request

	// Correlation returns the request correlation id.
	CorrelationID() string

	// User returns the user
	User() User

	// Logger returns the hexa logger customized for specific request.
	Logger() Logger

	// Translator returns the translator localized relative to the users request.
	Translator() Translator

	// ToMap method convert context to map to export and import.
	ToMap() Map
}

type defaultContext struct {
	context.Context
	locale        string
	request       *http.Request
	correlationID string
	user          User
	logger        Logger
	translator    Translator
}

// exportedCtx use export and import Context.
type exportedCtx struct {
	Locale        string      `json:"locale"`
	CorrelationID string      `json:"correlation_id"`
	UserId        interface{} `json:"user_id"`
}

func (e exportedCtx) validate() error {
	if e.CorrelationID == "" || e.UserId == "" {
		return tracer.Trace(errors.New("exported data is invalid"))
	}

	return nil
}

func (c defaultContext) Request() *http.Request {
	return c.request
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
		CorrelationID: c.CorrelationID(),
		UserId:        c.User().Identifier().Val(),
	})
}

// NewCtx returns new hexa Context.
// locale syntax is just same as HTTP Accept-Language header.
func NewCtx(request *http.Request, correlationID string, locale string, user User, logger Logger, translator Translator) Context {
	logger = tuneCtxLogger(request, correlationID, user, logger)
	translator = tuneCtxTranslator(locale, translator)

	return &defaultContext{
		Context:       context.Background(),
		locale:        locale,
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

	return NewCtx(nil, e.CorrelationID, e.Locale, u, l, t), nil
}

// tuneLogger function tune the logger for each context.
func tuneCtxLogger(r *http.Request, correlationID string, u User, logger Logger) Logger {
	tags := map[string]interface{}{
		"__guest__":          u.IsGuest(),
		"__user_id__":        u.Identifier().String(),
		"__username__":       u.GetUsername(),
		"__correlation_id__": correlationID,
	}

	if !u.IsGuest() {
		tags["__email__"] = u.GetEmail()
		tags["__phone__"] = u.GetPhone()
	}

	if r != nil {
		rid := r.Header.Get("X-Request-ID")
		if rid != "" {
			tags["__request_id__"] = rid
		}

		ip := gutil.IP(r)
		if ip != "" {
			tags["__ip__"] = ip
		}
	}

	logger = logger.WithFields(gutil.MapToKeyValue(tags)...)

	return logger
}

// tuneCtxTranslator localize translator for each context.
func tuneCtxTranslator(locale string, t Translator) Translator {
	if locale != "" {
		return t.Localize(locale)
	}

	return t.Localize()
}

// Assert defaultContext implements the hexa Context.
var _ Context = &defaultContext{}

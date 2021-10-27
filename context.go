package hexa

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

const (
	ContextKeyRequest       = "_ctx_request"        // value must be *http.Request (optional)
	ContextKeyCorrelationID = "_ctx_correlation_id" // value Must be cid string
	ContextKeyLocale        = "_ctx_local"          // value must be locale string (can be empty string)
	ContextKeyUser          = "_ctx_user"           // Value must be user
	ContextKeyLogger        = "_ctx_logger"         // value must be Logger interface
	ContextKeyTranslator    = "_ctx_translator"     // value must be Translator interface
	ContextKeyStore         = "_ctx_store"          // value must be Store interface.
)

var requiredContextKeys = []string{
	ContextKeyUser,
	ContextKeyCorrelationID,
	ContextKeyLocale,
	ContextKeyLogger,
	ContextKeyTranslator,
}

// Context is the hexa context to use in services.
// This context must be just a wrapper of the golang context.
type Context interface {
	context.Context

	// Request returns the current request and can be nil for not http requests.
	Request() *http.Request

	// CorrelationID returns the request's correlation id.
	CorrelationID() string

	// User returns the user
	User() User

	// Logger returns the hexa logger customized for specific request.
	Logger() Logger

	// Translator returns the translator localized relative to the users request.
	Translator() Translator

	// Store returns the context embedded store. its default
	// implementation is just a simple concurrency-safe map.
	Store() Store
}

type contextImpl struct {
	context.Context
}

func (c contextImpl) Store() Store {
	return c.Value(ContextKeyStore).(Store)
}

func (c contextImpl) Request() *http.Request {
	return c.Value(ContextKeyRequest).(*http.Request)
}

func (c contextImpl) CorrelationID() string {
	return c.Value(ContextKeyCorrelationID).(string)
}

func (c contextImpl) User() User {
	return c.Value(ContextKeyUser).(User)
}

func (c contextImpl) Logger() Logger {
	field := StringField
	u := c.User()
	r := c.Request()

	tags := []LogField{
		field("_user_type", string(u.Type())),
		field("_user_id", u.Identifier()),
		field("_username", u.Username()),
		field("_correlation_id", c.CorrelationID()),
	}

	if u.Type() == UserTypeRegular {
		tags = append(tags, field("_email", u.Email()))
		tags = append(tags, field("_phone", u.Phone()))
	}

	if r != nil {
		rid := r.Header.Get("X-Request-ID")
		if rid != "" {
			tags = append(tags, field("_request_id", rid))
		}

		if ip, port, err := net.SplitHostPort(gutil.IP(r)); err == nil {
			tags = append(tags, field("_ip", ip))
			tags = append(tags, field("_port", port))
		}
	}
	logger := c.Value(ContextKeyLogger).(Logger)

	return logger.With(tags...)
}

func (c contextImpl) Translator() Translator {
	t := c.Value(ContextKeyTranslator).(Translator)
	locale := c.Value(ContextKeyLocale).(string)

	if locale != "" {
		return t.Localize(locale)
	}

	return t.Localize()
}

func WithUser(c Context, user User) Context {
	return WithValue(c, ContextKeyUser, user)
}

// WithValue is just same as context.WithValue function but returns
// hexa Context.
func WithValue(c Context, key interface{}, value interface{}) Context {
	return MustNewContextFromRawContext(context.WithValue(c, key, value))
}

func NewContextFromRawContext(c context.Context) (Context, error) {
	// If the provided context is a hexa context, we don't need to create a new one.
	if hc, ok := c.(Context); ok {
		return hc, nil
	}

	if err := validateRawContext(c); err != nil {
		return nil, tracer.Trace(err)
	}

	return &contextImpl{Context: c}, nil
}

func MustNewContextFromRawContext(c context.Context) Context {
	return gutil.Must(NewContextFromRawContext(c)).(Context)
}

type ContextParams struct {
	Request       *http.Request
	CorrelationId string
	Locale        string
	User          User
	Logger        Logger
	Translator    Translator
	Store         Store // Optional
}

// NewContext returns new hexa Context.
// locale syntax is just same as HTTP Accept-Language header.
func NewContext(ctx context.Context, p ContextParams) Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if p.Store == nil {
		p.Store = newStore()
	}

	c := contextWithParams(ctx, p)
	return MustNewContextFromRawContext(c)
}

func contextWithParams(c context.Context, p ContextParams) context.Context {
	c = context.WithValue(c, ContextKeyRequest, p.Request)
	c = context.WithValue(c, ContextKeyCorrelationID, p.CorrelationId)
	c = context.WithValue(c, ContextKeyLocale, p.Locale)
	c = context.WithValue(c, ContextKeyUser, p.User)
	c = context.WithValue(c, ContextKeyLogger, p.Logger)
	c = context.WithValue(c, ContextKeyTranslator, p.Translator)
	c = context.WithValue(c, ContextKeyStore, p.Store)

	return c
}

// validateRawContext validate check whether a raw context
// can be converted to a hexa context or not.
func validateRawContext(c context.Context) error {
	if k := getMissedKeyInContext(c, requiredContextKeys...); k != "" {
		errMsg := fmt.Sprintf("can not find key %s in context keys to generate hexa context", k)
		return tracer.Trace(errors.New(errMsg))
	}

	// assert user type:
	if _, ok := c.Value(ContextKeyUser).(User); !ok {
		return tracer.Trace(errors.New("invalid user for hexa context"))
	}

	// Request must be *http.Request if exists:
	if v := c.Value(ContextKeyRequest); v != nil {
		if _, ok := v.(*http.Request); !ok {
			return tracer.Trace(errors.New("request type is invalid for hexa context"))
		}
	}

	// CorrelationId can not be empty
	if cid, ok := c.Value(ContextKeyCorrelationID).(string); !ok {
		return tracer.Trace(errors.New("correlation id type is invalid for hexa context"))
	} else if cid == "" {
		return tracer.Trace(errors.New("correlation id can not be empty for hexa context"))
	}

	// local must be string
	if _, ok := c.Value(ContextKeyLocale).(string); !ok {
		return tracer.Trace(errors.New("local type is invalid for hexa context"))
	}

	// assert logger type
	if _, ok := c.Value(ContextKeyLogger).(Logger); !ok {
		return tracer.Trace(errors.New("invalid logger for hexa context"))
	}

	// assert translator type:
	if _, ok := c.Value(ContextKeyTranslator).(Translator); !ok {
		return tracer.Trace(errors.New("invalid translator for hexa context"))
	}

	// assert store type:
	if _, ok := c.Value(ContextKeyStore).(Store); !ok {
		return tracer.Trace(errors.New("invalid store for hexa context"))
	}

	return nil
}

// Assert contextImpl implements the hexa Context.
var _ Context = &contextImpl{}

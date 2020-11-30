package hexa

import (
	"context"
	"errors"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"net"
	"net/http"
)

type (
	// Context is the hexa context to use in services.
	Context interface {
		context.Context

		// WithUser returns new instance of context with provided user.
		WithUser(User) Context

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
		ToMap(UserExporterImporter) (Map, error)
	}

	defaultContext struct {
		context.Context
		baseLogger     Logger
		baseTranslator Translator
		locale         string
		request        *http.Request
		correlationID  string
		user           User
		logger         Logger
		translator     Translator
	}

	// exportedCtx use export and import Context.
	exportedCtx struct {
		Locale        string `json:"locale"`
		CorrelationID string `json:"correlation_id"`
		User          Map    `json:"user"`
	}

	// ContextExporterImporter export and import the context
	ContextExporterImporter interface {
		Export(Context) (Map, error)
		Import(Map) (Context, error)
	}

	// contextExporterImporter export & import the context.
	contextExporterImporter struct {
		ue UserExporterImporter
		l  Logger
		t  Translator
	}
)

func (e exportedCtx) validate() error {
	if e.CorrelationID == "" || e.User == nil {
		return tracer.Trace(errors.New("exported data is invalid"))
	}

	return nil
}

func (c defaultContext) WithUser(user User) Context {
	return NewCtx(c.request, c.correlationID, c.locale, user, c.baseLogger, c.baseTranslator)
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

func (c defaultContext) ToMap(ue UserExporterImporter) (Map, error) {
	if ue == nil {
		return nil, tracer.Trace(errors.New("user exporter can not be nil"))
	}

	exportedUser, err := ue.Export(c.User())
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return gutil.StructToMap(exportedCtx{
		Locale:        c.locale,
		CorrelationID: c.CorrelationID(),
		User:          exportedUser,
	}), nil
}

func (ce *contextExporterImporter) Export(ctx Context) (Map, error) {
	return ctx.ToMap(ce.ue)
}

func (ce *contextExporterImporter) Import(m Map) (Context, error) {
	e := exportedCtx{}
	err := gutil.MapToStruct(m, &e)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	if err := e.validate(); err != nil {
		return nil, tracer.Trace(err)
	}

	u, err := ce.ue.Import(e.User)
	if err != nil {
		return nil, err
	}

	return NewCtx(nil, e.CorrelationID, e.Locale, u, ce.l, ce.t), nil
}

// NewCtx returns new hexa Context.
// locale syntax is just same as HTTP Accept-Language header.
func NewCtx(request *http.Request, correlationID string, locale string, user User, logger Logger, translator Translator) Context {
	ctx := &defaultContext{
		Context:        context.Background(),
		baseLogger:     logger,
		baseTranslator: translator,
		request:        request,
		locale:         locale,
		correlationID:  correlationID,
		user:           user,
		logger:         logger,
		translator:     tuneCtxTranslator(locale, translator),
	}

	// Bind context to the context's logger.
	ctx.logger = tuneCtxLogger(request, correlationID, user, logger).With(ctx)
	return ctx
}

// NewCtxExporterImporter returns new instance of the ContextExporterImporter to export and import context.
func NewCtxExporterImporter(ue UserExporterImporter, l Logger, t Translator) ContextExporterImporter {
	return &contextExporterImporter{
		ue: ue,
		l:  l,
		t:  t,
	}
}

// tuneLogger function tune the logger for each context.
func tuneCtxLogger(r *http.Request, correlationID string, u User, logger Logger) Logger {
	field := StringField

	tags := []LogField{
		field("__user_type__", string(u.Type())),
		field("__user_id__", u.Identifier().String()),
		field("__username__", u.Username()),
		field("__correlation_id__", correlationID),
	}

	if u.Type() == UserTypeRegular {
		tags = append(tags, field("__email__", u.Email()))
		tags = append(tags, field("__phone__", u.Phone()))
	}

	if r != nil {
		rid := r.Header.Get("X-Request-ID")
		if rid != "" {
			tags = append(tags, field("__request_id__", rid))
		}

		if ip, port, err := net.SplitHostPort(gutil.IP(r)); err == nil {
			tags = append(tags, field("__ip__", ip))
			tags = append(tags, field("__port__", port))
		}
	}

	logger = logger.WithFields(tags...)
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
var _ ContextExporterImporter = &contextExporterImporter{}

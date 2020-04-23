package hexalogger

import (
	"errors"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/getsentry/sentry-go"
	"net"
	"net/http"
)

type sentryLogger struct {
	hub *sentry.Hub
}

var (
	LogConfigKeySentryDSN          = "log.sentry.dsn"
	LogConfigKeySentryEnvirontment = "log.sentry.environment"
)

func (l *sentryLogger) Core() interface{} {
	return l.hub
}

func (l *sentryLogger) argsToMap(args ...interface{}) map[string]interface{} {
	// if key values is not odd, add another item to make it odd.
	if len(args)%2 != 0 {
		args = append(args, errMissingValue)
	}
	fields, _ := gutil.KeyValuesToMap(args...)
	return fields
}

func (l *sentryLogger) addArgsToScope(scope *sentry.Scope, args []interface{}) {
	if len(args) == 0 {
		return
	}
	fields := l.argsToMap(args...)
	for key, val := range fields {
		// Just keys that begin and end with "_", set as tags.
		if len(key) >= 2 && key[0] == '_' && key[len(key)-1] == '_' {
			scope.SetTag(key, fmt.Sprintf("%v", val))
		} else {
			scope.SetExtra(key, val)
		}
	}
}

func (l *sentryLogger) setUser(scope *sentry.Scope, user hexa.User, r *http.Request) {
	u := sentry.User{
		IPAddress: gutil.IP(r),
		Email:     user.Email(),
		ID:        user.Identifier().String(),
		Username:  user.Username(),
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		u.IPAddress = ip
	}

	scope.SetUser(u)
}

func (l *sentryLogger) setRequest(scope *sentry.Scope, r *http.Request) {

	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = fmt.Sprintf("%v", v)
	}
	var env map[string]string
	if addr, port, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		env = map[string]string{"REMOTE_ADDR": addr, "REMOTE_PORT": port}
	}

	scope.SetRequest(sentry.Request{
		URL:         r.URL.String(),
		Method:      r.Method,
		Data:        "",
		QueryString: r.URL.RawQuery,
		Cookies:     "",
		Headers:     headers,
		Env:         env,
	})
}

func (l *sentryLogger) With(ctx hexa.Context, args ...interface{}) hexa.Logger {
	hub := l.hub.Clone()
	scope := hub.Scope()

	r := ctx.Request()
	if r != nil {
		l.setRequest(scope, r)
	}

	l.setUser(scope, ctx.User(), r)

	l.addArgsToScope(scope, args)
	return NewSentryDriverWith(hub)
}

// WithFields get some fields and set check if field's key start and end
// with single '_' character, then insert it as tag, otherwise
// insert it as extra data.
func (l *sentryLogger) WithFields(args ...interface{}) hexa.Logger {
	hub := l.hub.Clone()
	l.addArgsToScope(hub.Scope(), args)
	return NewSentryDriverWith(hub)
}

func (l *sentryLogger) Debug(i ...interface{}) {
	// For now we do not capture debug messages in sentry.
}

func (l *sentryLogger) Info(i ...interface{}) {
	// For now we do not capture messages in info .
}

func (l *sentryLogger) Message(i ...interface{}) {
	l.hub.CaptureMessage(fmt.Sprint(i...))
}

func (l *sentryLogger) Warn(i ...interface{}) {
	// For now we do not capture message in warn.
}

func (l *sentryLogger) Error(i ...interface{}) {
	l.hub.CaptureException(errors.New(fmt.Sprint(i...)))
}

// NewSentryDriver return new instance of hexa logger with sentry driver.
func NewSentryDriver(config hexa.Config) (hexa.Logger, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         config.GetString(LogConfigKeySentryDSN),
		Debug:       config.GetBool("DEBUG"),
		Environment: config.GetString(LogConfigKeySentryEnvirontment),
	})
	if err != nil {
		return nil, err
	}
	return NewSentryDriverWith(sentry.NewHub(client, sentry.NewScope())), nil
}

// NewSentryDriverWith get the sentry hub and returns new instance
//of sentry driver for hexa logger.
func NewSentryDriverWith(hub *sentry.Hub) hexa.Logger {
	return &sentryLogger{hub}
}

// Assert sentryLogger implements hexa Logger.
var _ hexa.Logger = &sentryLogger{}

package hexalogger

import (
	"errors"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/getsentry/sentry-go"
)

type sentryLogger struct {
	hub *sentry.Hub
}

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

func (l *sentryLogger) With(ctx hexa.Context, args ...interface{}) hexa.Logger {
	hub := l.hub.Clone()
	scope := hub.Scope()

	user := ctx.User()

	// Set the user:
	scope.SetUser(sentry.User{
		Email:     "",
		ID:        user.Identifier().String(),
		IPAddress: "",
		Username:  user.GetUsername(),
	})

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

func (l *sentryLogger) Error(i ...interface{}) {
	l.hub.CaptureException(errors.New(fmt.Sprint(i...)))
}

// NewSentryDriver return new instance of hexa logger with sentry driver.
func NewSentryDriver(config hexa.Config) (hexa.Logger, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         config.GetString("SENTRY_DSN"),
		Debug:       config.GetBool("DEBUG"),
		Environment: config.GetString("SENTRY_ENVIRONMENT"),
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

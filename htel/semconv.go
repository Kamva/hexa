package htel

import "go.opentelemetry.io/otel/attribute"

//--------------------------------
// Extra key names for tracing
//--------------------------------

const (
	EnduserUsernameKey = attribute.Key("enduser.username")
	CorrelationIDKey   = attribute.Key("ctx.correlation_id")
)

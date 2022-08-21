package hexa

import (
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

// ErrFields checks if the provided error is a Hexa error, returns
// hexa error fields, otherwise returns regular error fields.
func ErrFields(err error) []hlog.Field {
	if hexaErrFields := hexaErrFields(err); len(hexaErrFields) != 0 {
		return hexaErrFields
	}

	return []hlog.Field{
		hlog.Err(err),
		hlog.ErrStack(tracer.Trace(err)),
	}
}

func hexaErrFields(err error) []hlog.Field {
	e := AsHexaErr(err)
	if e == nil {
		return nil
	}

	// Hexa error fields:
	fields := []hlog.Field{
		hlog.String("_error_id", e.ID()),
		hlog.Int("_http_status", e.HTTPStatus()),
	}
	for k, v := range e.Data() {
		fields = append(fields, hlog.Any(k, v))
	}
	for k, v := range e.ReportData() {
		fields = append(fields, hlog.Any(k, v))
	}

	// If exists error and error is traced,print its stack.
	fields = append(fields, hlog.ErrStack(tracer.MoveStackIfNeeded(e, e.InternalError())))
	if e.InternalError() != nil {
		fields = append(fields, hlog.Err(e.InternalError()))
	}

	return fields
}

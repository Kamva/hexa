package kitty

import (
	"encoding/json"
)

// HttpRespBody is the http response body format
type HttpRespBody struct {
	debug     bool
	debugData Map

	Code    string `json:"code" mapstructure:"code"`
	Message string `json:"message" mapstructure:"message"`
	Data    Map   `json:"data" mapstructure:"data"`
}

// MarshalJSON marshall the body to json value.
func (b HttpRespBody) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"code":    b.Code,
		"message": b.Message,
		"data":    b.Data,
	}

	if b.debug {
		m["__debug__"] = b.debugData
	}

	return json.Marshal(m)
}

// Debug set debug flag and debug data.
func (b HttpRespBody) Debug(debug bool, debugData Map) HttpRespBody {
	b.debug = debug
	b.debugData = debugData

	return b
}

// NewBody return new instance of the HttpRespBody
func NewBody(code string, msg string, data Map) HttpRespBody {
	return HttpRespBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

// Assert HttpRespBody implements the json unmarshaller.
var _ json.Marshaler = &HttpRespBody{}

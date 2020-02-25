package kitty

import (
	"encoding/json"
	"github.com/Kamva/gutil"
)

// Data is the response body data.
type Data map[string]interface{}

// HttpRespBody is the http response body format
type HttpRespBody struct {
	debug     bool
	debugData Data

	Code    string `json:"code" mapstructure:"code"`
	Message string `json:"message" mapstructure:"message"`
	Data    Data   `json:"data" mapstructure:"data"`
}

// MarshalJSON marshall the body to json value.
func (b HttpRespBody) MarshalJSON() ([]byte, error) {
	m := gutil.StructToMap(b)

	if b.debug {
		m["__debug__"] = b.debugData
	}

	return json.Marshal(m)
}

// Debug set debug flag and debug data.
func (b HttpRespBody) Debug(debug bool, debugData Data) {
	b.debug = debug
	b.debugData = debugData
}

// NewBody return new instance of the HttpRespBody
func NewBody(code string, msg string, data Data) HttpRespBody {
	return HttpRespBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

// Assert HttpRespBody implements the json unmarshaller.
var _ json.Marshaler = &HttpRespBody{}

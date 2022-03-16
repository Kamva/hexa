//go:generate easyjson

package hexa

// HttpRespBody is the http response body format
//easyjson:json
type HttpRespBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
	Debug   interface{}  `json:"debug,omitempty"` // Set this value to nil when you are on production mode.
}

// SetDebug set debug flag and debug data.
func (b HttpRespBody) SetDebug(debugData interface{}) HttpRespBody {
	b.Debug = debugData
	return b
}

// NewBody return new instance of the HttpRespBody
func NewBody(code string, msg string, data interface{}) HttpRespBody {
	return HttpRespBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

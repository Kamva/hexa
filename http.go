//go:generate easyjson

package hexa

// HttpRespBody is the http response body format
//easyjson:json
type HttpRespBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Debug   interface{} `json:"debug,omitempty"` // Set this value to nil when you are on production mode.
}

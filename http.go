//go:generate easyjson

package hexa

// HTTPRespBody is the http response body format
// easyjson:json
type HTTPRespBody struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Debug   any    `json:"debug,omitempty"` // Set this value to nil when you are on production mode.
}

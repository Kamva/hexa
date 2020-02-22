package kitty

// Data is the response body data.
type Data map[string]interface{}

// HttpRespBody is the http response body format
type HttpRespBody struct {
	Code    string `json:"code" mapstructure:"code"`
	Message string `json:"message" mapstructure:"message"`
	Data    Data   `json:"data" mapstructure:"data"`
}

// NewBody return new instance of the HttpRespBody
func NewBody(code string, msg string, data Data) HttpRespBody {
	return HttpRespBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

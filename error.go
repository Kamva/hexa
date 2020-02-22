package kitty

// ErrorData is extra data of the error to show to the user.
type ErrorData map[string]interface{}

// ErrorParams is parameters of the error to use in translation,...
type ErrorParams map[string]interface{}

type Error interface {
	error

	// Code return the error edentifier code
	Code() string

	// Key returns unique key for each error to use as translation key,...
	Key() string

	// Params returns params of the error to use in translation,...
	Params() ErrorParams

	// SetParams set the error parameters to use in error translation,...
	SetParams(params ErrorParams) Error

	// Data returns the extra data of the error.
	Data() ErrorData

	// SetData set the error data as extra data of the error to show to the user.
	SetData(data ErrorData) Error
}

// defaultError implements Error interface.
type defaultError struct {
	code   string
	key    string
	err    string
	params ErrorParams
	data   ErrorData
}

func (e defaultError) Error() string {
	return e.err
}

func (e defaultError) Code() string {
	return e.code
}

func (e defaultError) Key() string {
	return e.code
}

func (e defaultError) Params() ErrorParams {
	return e.params
}

func (e defaultError) SetParams(params ErrorParams) Error {
	e.params = params
	return e
}

func (e defaultError) Data() ErrorData {
	return e.data
}

func (e defaultError) SetData(data ErrorData) Error {
	e.data = data
	return e
}

// NewError returns new Error instance.
func NewError(code string, key string, err string) Error {
	return defaultError{
		code: code,
		key:  key,
		err:  err,
	}
}

// Assert defaultError implements the Error interface.
var _ Error = defaultError{}

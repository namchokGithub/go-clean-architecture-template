package common

// ErrorResponse is the standard error response envelope.
type ErrorResponse struct {
	IsError   bool   `json:"is_error"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	ErrorType string `json:"error_type"`
	Errors    any    `json:"errors"`
}

// Error is the handler-layer error type returned from handler functions.
// Echo's global error handler maps it to an HTTP response.
type Error struct {
	HTTPStatus int
	Code       string
	Message    string
	ErrorType  string
	Errors     any
}

func (e *Error) Error() string { return e.Message }

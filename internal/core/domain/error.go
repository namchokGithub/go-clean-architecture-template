package domain

import "fmt"

// ErrorCode maps to HTTP status via its first 3 digits (e.g., "4001" → HTTP 400).
type ErrorCode string

const (
	ErrorCodeBadRequest   ErrorCode = "4001"
	ErrorCodeUnauthorized ErrorCode = "4011"
	ErrorCodeForbidden    ErrorCode = "4031"
	ErrorCodeNotFound     ErrorCode = "4041"
	ErrorCodeConflict     ErrorCode = "4091"
	ErrorCodeInternal     ErrorCode = "5001"
)

// SubError represents a single field-level validation error.
type SubError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error is the domain-layer error type. Handlers map it to HTTP responses.
type Error struct {
	Code    ErrorCode
	Message string
	Errors  []SubError
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a domain error with field-level sub-errors.
func NewError(code ErrorCode, errs []SubError) *Error {
	return &Error{Code: code, Errors: errs}
}

// NewErrorString creates a domain error with a plain message.
func NewErrorString(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message}
}

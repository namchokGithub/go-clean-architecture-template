package common

import (
	"net/http"

	"gitlab.company.com/projectname/internal/handler/base"
)

var appPrefix = "APP"

// SetAppPrefix sets the prefix used in all response codes (e.g., "SDD").
// Call once during protocol initialization.
func SetAppPrefix(prefix string) {
	appPrefix = prefix
}

// GetAppPrefix returns the current app prefix.
func GetAppPrefix() string {
	return appPrefix
}

func code(suffix string) string {
	return appPrefix + "-" + suffix
}

// SubError is used in validation error responses.
type SubError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Success responses

func NewSuccessResponse(data any) base.DefaultResponse {
	return base.DefaultResponse{IsError: false, Code: code("200"), Message: "Success", Data: data}
}

func NewCreatedResponse(data any) base.DefaultResponse {
	return base.DefaultResponse{IsError: false, Code: code("201"), Message: "Created", Data: data}
}

// Error constructors — each returns *Error which Echo's error handler renders.

func NewBadRequest(errs any) *Error {
	return &Error{HTTPStatus: http.StatusBadRequest, Code: code("4001"), Message: "Bad Request", ErrorType: "BAD_REQUEST", Errors: errs}
}

func NewUnAuthorized() *Error {
	return &Error{HTTPStatus: http.StatusUnauthorized, Code: code("4011"), Message: "Unauthorized", ErrorType: "UNAUTHORIZED"}
}

func NewForbidden() *Error {
	return &Error{HTTPStatus: http.StatusForbidden, Code: code("4031"), Message: "Forbidden", ErrorType: "FORBIDDEN"}
}

func NewNotFound() *Error {
	return &Error{HTTPStatus: http.StatusNotFound, Code: code("4041"), Message: "Not Found", ErrorType: "NOT_FOUND"}
}

func NewConflict() *Error {
	return &Error{HTTPStatus: http.StatusConflict, Code: code("4091"), Message: "Conflict", ErrorType: "CONFLICT"}
}

func NewInternalServerError() *Error {
	return &Error{HTTPStatus: http.StatusInternalServerError, Code: code("5001"), Message: "Internal Server Error", ErrorType: "INTERNAL_SERVER_ERROR"}
}

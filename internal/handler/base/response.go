package base

// DefaultResponse is the standard success response envelope.
type DefaultResponse struct {
	IsError bool   `json:"is_error"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

package domain

// ContextKey is the typed key for context values to prevent collisions.
type ContextKey string

const (
	Lang            ContextKey = "Lang"
	RequestID       ContextKey = "RequestID"
	OS              ContextKey = "OS"
	Browser         ContextKey = "Browser"
	USER_ID         ContextKey = "USER_ID"
	USER_TYPE       ContextKey = "USER_TYPE"
	ROLE_ID         ContextKey = "ROLE_ID"
	TOKEN           ContextKey = "TOKEN"
	EMAIL           ContextKey = "EMAIL"
	USER            ContextKey = "USER"
	ROLE_PERMISSION ContextKey = "ROLE_PERMISSION"
	MENU_ID         ContextKey = "MENU_ID"
)

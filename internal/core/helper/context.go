package helper

import (
	"context"

	"gitlab.company.com/projectname/internal/core/domain"
)

// GetUserID extracts USER_ID from context. Returns 0 if missing.
func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(domain.USER_ID).(int64)
	return v
}

// GetUserType extracts USER_TYPE from context.
func GetUserType(ctx context.Context) string {
	v, _ := ctx.Value(domain.USER_TYPE).(string)
	return v
}

// GetRoleID extracts ROLE_ID from context.
func GetRoleID(ctx context.Context) int64 {
	v, _ := ctx.Value(domain.ROLE_ID).(int64)
	return v
}

// GetLang extracts Lang from context. Returns "th" as default.
func GetLang(ctx context.Context) string {
	v, _ := ctx.Value(domain.Lang).(string)
	if v == "" {
		return "th"
	}
	return v
}

// GetToken extracts the raw JWT token string from context.
func GetToken(ctx context.Context) string {
	v, _ := ctx.Value(domain.TOKEN).(string)
	return v
}

package jwt

import (
	"context"
	"fmt"
	"strings"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// contextKey is the typed key used by this package for JWT values in context.
type contextKey string

const (
	CtxUserID   contextKey = "jwt_user_id"
	CtxUserType contextKey = "jwt_user_type"
	CtxRoleID   contextKey = "jwt_role_id"
	CtxEmail    contextKey = "jwt_email"
	CtxToken    contextKey = "jwt_token"
)

// Claims is the JWT payload.
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserType string `json:"user_type"`
	RoleID   int64  `json:"role_id"`
	Email    string `json:"email"`
	gojwt.RegisteredClaims
}

// GetUserID extracts UserID from context (set by AuthMiddleware).
func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(CtxUserID).(int64)
	return v
}

// GetUserType extracts UserType from context.
func GetUserType(ctx context.Context) string {
	v, _ := ctx.Value(CtxUserType).(string)
	return v
}

// GetRoleID extracts RoleID from context.
func GetRoleID(ctx context.Context) int64 {
	v, _ := ctx.Value(CtxRoleID).(int64)
	return v
}

// GetEmail extracts Email from context.
func GetEmail(ctx context.Context) string {
	v, _ := ctx.Value(CtxEmail).(string)
	return v
}

// GetToken extracts the raw token string from context.
func GetToken(ctx context.Context) string {
	v, _ := ctx.Value(CtxToken).(string)
	return v
}

// AuthMiddleware validates the Bearer JWT and injects claims into the request context.
func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return echo.ErrUnauthorized
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			claims := &Claims{}
			token, err := gojwt.ParseWithClaims(tokenStr, claims, func(t *gojwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return echo.ErrUnauthorized
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxUserType, claims.UserType)
			ctx = context.WithValue(ctx, CtxRoleID, claims.RoleID)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			ctx = context.WithValue(ctx, CtxToken, tokenStr)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

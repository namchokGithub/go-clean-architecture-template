package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/core/constant"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/handler/common"
)

// APIKeyGuard rejects requests that do not supply the correct X-Api-Key header.
func APIKeyGuard(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get(constant.HeaderAPIKey) != expectedKey {
				return common.NewUnAuthorized()
			}
			return next(c)
		}
	}
}

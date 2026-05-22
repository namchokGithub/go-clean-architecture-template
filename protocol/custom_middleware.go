package protocol

import (
	"context"

	"github.com/labstack/echo/v4"
	"gitlab.company.com/projectname/internal/core/constant"
	"gitlab.company.com/projectname/internal/core/domain"
	"gitlab.company.com/projectname/internal/handler/common"
	pkgjwt "gitlab.company.com/projectname/pkg/jwt"
	"gitlab.company.com/projectname/pkg/logx"
)

// globalContextMiddleware injects Lang, OS, Browser into the request context.
func (a *App) globalContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			lang := c.Request().Header.Get(constant.HeaderAcceptLanguage)
			if lang == "" {
				lang = "th"
			}
			ctx = context.WithValue(ctx, domain.Lang, lang)
			ctx = context.WithValue(ctx, domain.OS, c.Request().Header.Get("X-OS"))
			ctx = context.WithValue(ctx, domain.Browser, c.Request().Header.Get("User-Agent"))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// afterAuthMiddleware promotes JWT claims from pkg/jwt context keys to domain context keys.
// Also loads the full user + role permissions (Redis cache → DB fallback — TODO).
func (a *App) afterAuthMiddleware() echo.MiddlewareFunc {
	log := logx.GetLog()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			token := pkgjwt.GetToken(ctx)
			if token == "" {
				return common.NewUnAuthorized()
			}

			// Promote JWT-layer claims to domain context keys.
			ctx = context.WithValue(ctx, domain.USER_ID, pkgjwt.GetUserID(ctx))
			ctx = context.WithValue(ctx, domain.USER_TYPE, pkgjwt.GetUserType(ctx))
			ctx = context.WithValue(ctx, domain.ROLE_ID, pkgjwt.GetRoleID(ctx))
			ctx = context.WithValue(ctx, domain.EMAIL, pkgjwt.GetEmail(ctx))
			ctx = context.WithValue(ctx, domain.TOKEN, token)

			// TODO: load user from Redis cache ({token}_user)
			// TODO: fall back to DB on cache miss and re-cache
			// TODO: load role permissions from Redis ({token}_rolePermission)
			// TODO: inject domain.USER and domain.ROLE_PERMISSION into ctx

			log.WithField("token_prefix", safeTokenPrefix(token)).Debug("afterAuth")
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func safeTokenPrefix(s string) string {
	if len(s) > 8 {
		return s[:8] + "..."
	}
	return s
}

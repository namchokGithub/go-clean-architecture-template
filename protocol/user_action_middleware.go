package protocol

import (
	"github.com/labstack/echo/v4"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/core/domain"
	"gitlab.socket9.com/sritrang/super-driver/backend/pkg/logx"
)

// userActionMiddleware logs every authenticated request.
func (a *App) userActionMiddleware() echo.MiddlewareFunc {
	log := logx.GetLog()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			userID, _ := ctx.Value(domain.USER_ID).(int64)
			log.WithFields(map[string]interface{}{
				"user_id": userID,
				"method":  c.Request().Method,
				"path":    c.Path(),
			}).Info("user_action")
			// TODO: persist to tbl_user_logs via logrus DB hook
			return next(c)
		}
	}
}

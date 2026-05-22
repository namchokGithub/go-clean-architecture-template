package protocol

import "github.com/labstack/echo/v4"

// permissionMiddleware checks ROLE_PERMISSION[menuID] for required access.
// Currently a pass-through; populate when RBAC is fully wired.
func (a *App) permissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: extract MENU_ID from route, check ROLE_PERMISSION from context
			return next(c)
		}
	}
}

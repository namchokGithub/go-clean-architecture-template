package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// EnforceJSONContentType normalizes Content-Type for POST/PUT/PATCH requests.
func EnforceJSONContentType() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method == "POST" || method == "PUT" || method == "PATCH" {
				ct := c.Request().Header.Get("Content-Type")
				if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "multipart/form-data") {
					c.Request().Header.Set("Content-Type", "application/json")
				}
			}
			return next(c)
		}
	}
}

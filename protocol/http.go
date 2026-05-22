package protocol

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gitlab.company.com/projectname/internal/core/constant"
	"gitlab.company.com/projectname/internal/handler/middleware"
	"gitlab.company.com/projectname/internal/handler/validator"
	pkgjwt "gitlab.company.com/projectname/pkg/jwt"
)

func (a *App) newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = middleware.ErrorHandler()
	validator.Register(e)

	// Global middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, constant.HeaderDeviceModel, constant.HeaderPlatform, constant.HeaderAcceptLanguage, constant.HeaderRequestID, constant.HeaderRequestChannel, constant.HeaderAppID, constant.HeaderAPIKey, constant.HeaderRequestClient, constant.HeaderDeviceID},
		AllowCredentials: true,
	}))
	e.Use(a.globalContextMiddleware())

	// Auth-gated group
	authGroup := e.Group("")
	authGroup.Use(pkgjwt.AuthMiddleware(a.conf.Key.JWTSecret))
	authGroup.Use(a.afterAuthMiddleware())
	authGroup.Use(a.permissionMiddleware())
	authGroup.Use(a.userActionMiddleware())

	// Register routes (public + auth-gated)
	a.handler.RegisterRoutes(e, authGroup)

	return e
}

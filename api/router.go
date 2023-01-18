package api

import (
	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/api/middleware"
	"github.com/gizmo-ds/misstodon/api/nodeinfo"
	v1 "github.com/gizmo-ds/misstodon/api/v1"
	"github.com/gizmo-ds/misstodon/api/wellknown"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Router(e *echo.Echo) {
	e.HTTPErrorHandler = httperror.ErrorHandler
	e.Use(middleware.SetContextData)
	if global.Config.Logger.RequestLogger {
		e.Use(middleware.Logger)
	}

	wellknownApi := e.Group("/.well-known", echomiddleware.CORS())
	wellknownApi.GET("/nodeinfo", wellknown.NodeInfo)

	e.GET("/nodeinfo/2.0", nodeinfo.NodeInfo, echomiddleware.CORS())

	v1Api := e.Group("/api/v1", echomiddleware.CORS())
	v1Api.GET("/instance", v1.Instance)
	v1Api.GET("/accounts/verify_credentials", v1.AccountsVerifyCredentials)
}

package api

import (
	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/api/middleware"
	"github.com/gizmo-ds/misstodon/api/nodeinfo"
	"github.com/gizmo-ds/misstodon/api/oauth"
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

	{
		wellknown.Router(e)
		nodeinfo.Router(e)
		v1Api := e.Group("/api/v1", echomiddleware.CORS())
		oauth.Router(e)
		v1.InstanceRouter(v1Api)
		v1.AccountsRouter(v1Api)
		v1.ApplicationRouter(v1Api)
	}
	{
		paramServerApi := e.Group("/:proxyServer")
		wellknown.Router(paramServerApi)
		nodeinfo.Router(paramServerApi)
		v1Api := paramServerApi.Group("/api/v1", echomiddleware.CORS())
		oauth.Router(paramServerApi)
		v1.InstanceRouter(v1Api)
		v1.AccountsRouter(v1Api)
		v1.ApplicationRouter(v1Api)
	}
}

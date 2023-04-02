package api

import (
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/api/nodeinfo"
	"github.com/gizmo-ds/misstodon/internal/api/oauth"
	"github.com/gizmo-ds/misstodon/internal/api/v1"
	"github.com/gizmo-ds/misstodon/internal/api/v2"
	"github.com/gizmo-ds/misstodon/internal/api/wellknown"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	e.HTTPErrorHandler = httperror.ErrorHandler
	e.Use(
		middleware.SetContextData,
		middleware.Recover())
	if global.Config.Logger.RequestLogger {
		e.Use(middleware.Logger)
	}
	for _, group := range []*echo.Group{
		e.Group(""),
		e.Group("/:proxyServer"),
	} {
		wellknown.Router(group)
		nodeinfo.Router(group)
		oauth.Router(group)
		v1Api := group.Group("/api/v1", middleware.CORS())
		group.GET("/static/missing.png", v1.MissingImageHandler)
		v1.InstanceRouter(v1Api)
		v1.AccountsRouter(v1Api)
		v1.ApplicationRouter(v1Api)
		v1.StatusesRouter(v1Api)
		v1.StreamingRouter(v1Api)
		v1.TimelinesRouter(v1Api)
		v1.TrendsRouter(v1Api)
		v1.MediaRouter(v1Api)
		v2.MediaRouter(v1Api)
	}
}

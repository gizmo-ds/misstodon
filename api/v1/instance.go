package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InstanceRouter(e *echo.Group) {
	group := e.Group("/instance", middleware.CORS())
	group.GET("", Instance)
	group.GET("/peers", InstancePeers)
}

func Instance(c echo.Context) error {
	info, err := misskey.Instance(
		c.Get("server").(string),
		global.AppVersion)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func InstancePeers(c echo.Context) error {
	peers, err := misskey.InstancePeers(c.Get("server").(string))
	if err != nil {
		return err
	}
	if len(peers) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}
	return c.JSON(http.StatusOK, peers)
}

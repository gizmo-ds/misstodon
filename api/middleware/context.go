package middleware

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/labstack/echo/v4"
)

func SetContextData(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		proxyServer, ok := utils.StrEvaluation(
			c.Param("proxyServer"),
			c.QueryParam("server"),
			c.Request().Header.Get("x-proxy-server"),
			global.Config.Proxy.FallbackServer)
		if !ok {
			if !utils.Contains([]string{
				"/.well-known/nodeinfo",
				"/nodeinfo/2.0",
			}, c.Path()) {
				return c.JSON(http.StatusBadRequest, httperror.ServerError{
					Error: "no proxy server specified",
				})
			}
		}
		c.Response().Header().Set("User-Agent", "misstodon/"+global.AppVersion)
		c.Response().Header().Set("X-Proxy-Server", proxyServer)
		c.Set("server", proxyServer)
		return next(c)
	}
}

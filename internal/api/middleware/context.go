package middleware

import (
	"net/http"
	"strings"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/labstack/echo/v4"
)

func SetContextData(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var hostProxyServer string
		host := c.Request().Host
		if strings.HasPrefix(host, "mt_") {
			tmp := strings.Split(host[3:], ".")[0]
			tmp = strings.ReplaceAll(tmp, "__", "+")
			arr := strings.Split(tmp, "_")
			if len(arr) > 1 {
				tmp = strings.Join(arr, ".")
				hostProxyServer = strings.ReplaceAll(tmp, "+", "_")
			}
		}
		proxyServer, ok := utils.StrEvaluation(
			hostProxyServer,
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
		c.Set("proxy-server", proxyServer)
		return next(c)
	}
}

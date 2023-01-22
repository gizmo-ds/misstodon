package wellknown

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router(e any) {
	var group *echo.Group
	switch e.(type) {
	case *echo.Echo:
		group = e.(*echo.Echo).Group("/.well-known", middleware.CORS())
	case *echo.Group:
		group = e.(*echo.Group).Group("/.well-known", middleware.CORS())
	}
	group.GET("/nodeinfo", NodeInfo)
	group.GET("/webfinger", WebFinger)
}

func NodeInfo(c echo.Context) error {
	server := c.Get("server").(string)
	href := "https://" + c.Request().Host + "/nodeinfo/2.0"
	if server != "" {
		href += "?server=" + server
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"links": []map[string]string{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				"href": href,
			},
		},
	})
}

func WebFinger(c echo.Context) error {
	resource := c.QueryParam("resource")
	if resource == "" {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: "resource is required",
		})
	}
	return misskey.WebFinger(c.Get("server").(string), resource, c.Response().Writer)
}

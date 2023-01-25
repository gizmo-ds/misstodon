package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func StatusesRouter(e *echo.Group) {
	group := e.Group("/statuses")
	group.GET("/:id", StatusHandler)
	group.GET("/:id/context", StatusContextHandler)
}

func StatusHandler(c echo.Context) error {
	server := c.Get("server").(string)
	id := c.Param("id")
	accessToken, _ := utils.GetHeaderToken(c.Request().Header)
	info, err := misskey.StatusSingle(server, accessToken, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func StatusContextHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"ancestors":   []interface{}{},
		"descendants": []interface{}{},
	})
}

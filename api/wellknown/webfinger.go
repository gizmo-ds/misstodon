package wellknown

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func WebFinger(c echo.Context) error {
	resource := c.QueryParam("resource")
	if resource == "" {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: "resource is required",
		})
	}
	return misskey.WebFinger(c.Get("server").(string), resource, c.Response().Writer)
}

package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func Instance(c echo.Context) error {
	info, err := misskey.Instance(c.Get("server").(string))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

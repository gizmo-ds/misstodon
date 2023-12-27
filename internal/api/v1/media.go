package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/labstack/echo/v4"
)

func MediaRouter(e *echo.Group) {
	group := e.Group("/media")
	group.POST("", MediaUploadHandler)
}

func MediaUploadHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	description := c.FormValue("description")

	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}

	ma, err := misskey.MediaUpload(ctx, file, description)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ma)
}

package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
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

	accessToken, err := utils.GetHeaderToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
	}
	server := c.Get("server").(string)
	ma, err := misskey.MediaUpload(server, accessToken, file, description)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ma)
}

package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func StatusesRouter(e *echo.Group) {
	group := e.Group("/statuses")
	group.GET("/:id", StatusHandler)
	group.GET("/:id/context", StatusContextHandler)
	group.POST("/:id/bookmark", StatusBookmark)
	group.POST("/:id/unbookmark", StatusUnBookmark)
}

func StatusHandler(c echo.Context) error {
	server := c.Get("server").(string)
	id := c.Param("id")
	token, _ := utils.GetHeaderToken(c.Request().Header)
	info, err := misskey.StatusSingle(server, token, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func StatusContextHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"ancestors":   []any{},
		"descendants": []any{},
	})
}

func StatusBookmark(c echo.Context) error {
	server := c.Get("server").(string)
	id := c.Param("id")
	token, err := utils.GetHeaderToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
	}

	status, err := misskey.StatusBookmark(server, token, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

func StatusUnBookmark(c echo.Context) error {
	server := c.Get("server").(string)
	id := c.Param("id")
	token, err := utils.GetHeaderToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
	}

	status, err := misskey.StatusUnBookmark(server, token, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

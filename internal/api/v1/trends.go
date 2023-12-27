package v1

import (
	"net/http"
	"strconv"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/labstack/echo/v4"
)

func TrendsRouter(e *echo.Group) {
	group := e.Group("/trends")
	group.GET("/tags", TrendsTagsHandler)
	group.GET("/statuses", TrendsStatusHandler)
}

func TrendsTagsHandler(c echo.Context) error {
	limit := 10
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = v
		if limit > 20 {
			limit = 20
		}
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	ctx, _ := misstodon.ContextWithEchoContext(c)

	tags, err := misskey.TrendsTags(ctx, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(tags))
}

func TrendsStatusHandler(c echo.Context) error {
	limit := 20
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = v
		if limit > 30 {
			limit = 30
		}
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	ctx, _ := misstodon.ContextWithEchoContext(c)
	statuses, err := misskey.TrendsStatus(ctx, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(statuses))
}

package v1

import (
	"net/http"
	"strconv"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
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
	server := c.Get("server").(string)
	accessToken, _ := utils.GetHeaderToken(c.Request().Header)

	tags, err := misskey.TrendsTags(server, accessToken, limit, offset)
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
	server := c.Get("server").(string)
	accessToken, _ := utils.GetHeaderToken(c.Request().Header)
	statuses, err := misskey.TrendsStatus(server, accessToken, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(statuses))
}

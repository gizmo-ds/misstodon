package v1

import (
	"net/http"
	"strconv"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/labstack/echo/v4"
)

func TimelinesRouter(e *echo.Group) {
	group := e.Group("/timelines")
	group.GET("/public", TimelinePublicHandler)
	group.GET("/home", TimelineHomeHandler)
	group.GET("/tag/:hashtag", TimelineHashtag)
}

func TimelinePublicHandler(c echo.Context) error {
	ctx, _ := misstodon.ContextWithEchoContext(c)
	limit := 20
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}
	timelineType := models.TimelinePublicTypeRemote
	if c.QueryParam("local") == "true" {
		timelineType = models.TimelinePublicTypeLocal
	}
	list, err := misskey.TimelinePublic(ctx,
		timelineType, c.QueryParam("only_media") == "true", limit,
		c.QueryParam("max_id"), c.QueryParam("min_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

func TimelineHomeHandler(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	limit := 20
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}
	list, err := misskey.TimelineHome(ctx,
		limit, c.QueryParam("max_id"), c.QueryParam("min_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

func TimelineHashtag(c echo.Context) error {
	ctx, _ := misstodon.ContextWithEchoContext(c)

	limit := 20
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}

	list, err := misskey.SearchStatusByHashtag(ctx, c.Param("hashtag"),
		limit, c.QueryParam("max_id"), c.QueryParam("since_id"), c.QueryParam("min_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

package v1

import (
	"net/http"
	"strconv"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func TimelinesRouter(e *echo.Group) {
	group := e.Group("/timelines")
	group.GET("/public", TimelinePublicHandler)
}

func TimelinePublicHandler(c echo.Context) error {
	server := c.Get("server").(string)
	accessToken, _ := utils.GetHeaderToken(c.Request().Header)
	limit := 20
	if _limit := c.QueryParam("limit"); _limit != "" {
		if v, err := strconv.Atoi(_limit); err != nil {
			limit = v
		}
	}
	timelineType := models.TimelinePublicTypeRemote
	if c.QueryParam("local") == "true" {
		timelineType = models.TimelinePublicTypeLocal
	}
	list, err := misskey.TimelinePublic(server, accessToken,
		timelineType, c.QueryParam("only_media") == "true", limit,
		c.QueryParam("max_id"), c.QueryParam("min_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}

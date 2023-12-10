package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func NotificationsRouter(e *echo.Group) {
	group := e.Group("/notifications")
	group.GET("", NotificationsHandler)
}

func NotificationsHandler(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	_ = ctx
	var query struct {
		MaxId   string `query:"max_id"`
		MinId   string `query:"min_id"`
		SinceId string `query:"since_id"`
		Limit   int    `query:"limit"`
	}
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}

	getTypes := func(name string) []models.NotificationType {
		types := lo.Map(c.QueryParams()[name], func(item string, _ int) models.NotificationType { return models.NotificationType(item) })
		types = lo.Filter(types, func(item models.NotificationType, _ int) bool {
			return item != "" && item.ToMkNotificationType() != models.MkNotificationTypeUnknown
		})
		return types
	}

	types := getTypes("types[]")
	excludeTypes := getTypes("exclude_types[]")

	result, err := misskey.NotificationsGet(ctx,
		query.Limit, query.SinceId, query.MinId, query.MaxId,
		types, excludeTypes, "")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

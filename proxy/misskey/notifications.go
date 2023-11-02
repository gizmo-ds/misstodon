package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func NotificationsGet(ctx Context,
	limit int, sinceId, minId, maxId string,
	types, excludeTypes []models.NotificationType, accountId string,
) ([]models.Notification, error) {
	limit = utils.NumRangeLimit(limit, 1, 100)

	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	_excludeTypes := lo.Map(excludeTypes,
		func(item models.NotificationType, _ int) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	_excludeTypes = append(_excludeTypes, models.MkNotificationTypeAchievementEarned)
	body["excludeTypes"] = _excludeTypes
	_includeTypes := lo.Map(types,
		func(item models.NotificationType, _ int) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	if len(_includeTypes) > 0 {
		body["includeTypes"] = _includeTypes
	}

	var result []models.MkNotification
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i/notifications"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	notifications := lo.Map(result, func(item models.MkNotification, _ int) models.Notification {
		n, err := item.ToNotification(ctx.ProxyServer())
		if err == nil {
			return n
		}
		return models.Notification{}
	})
	notifications = lo.Filter(notifications, func(item models.Notification, _ int) bool {
		return item.Type != models.NotificationTypeUnknown
	})
	return notifications, nil
}

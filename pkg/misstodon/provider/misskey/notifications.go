package misskey

import (
	"net/http"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/pkg/errors"
)

func NotificationsGet(ctx misstodon.Context,
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
	_excludeTypes := slice.Map(excludeTypes,
		func(_ int, item models.NotificationType) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	_excludeTypes = append(_excludeTypes, models.MkNotificationTypeAchievementEarned)
	if slice.Contain(_excludeTypes, models.MkNotificationTypeMention) {
		_excludeTypes = append(_excludeTypes, models.MkNotificationTypeReply)
	}
	body["excludeTypes"] = _excludeTypes
	_includeTypes := slice.Map(types,
		func(_ int, item models.NotificationType) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	if slice.Contain(_includeTypes, models.MkNotificationTypeMention) {
		_includeTypes = append(_includeTypes, models.MkNotificationTypeReply)
	}
	if len(_includeTypes) > 0 {
		body["includeTypes"] = _includeTypes
	}

	var result []models.MkNotification
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/i/notifications")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	notifications := slice.Map(result, func(_ int, item models.MkNotification) models.Notification {
		n, err := item.ToNotification(ctx)
		if err == nil {
			return n
		}
		return models.Notification{Type: models.NotificationTypeUnknown}
	})
	notifications = slice.Filter(notifications, func(_ int, item models.Notification) bool {
		return item.Type != models.NotificationTypeUnknown
	})
	return notifications, nil
}

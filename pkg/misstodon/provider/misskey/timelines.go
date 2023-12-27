package misskey

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/pkg/errors"
)

func TimelinePublic(ctx misstodon.Context,
	timelineType models.TimelinePublicType, onlyMedia bool,
	limit int, maxId, minId string) ([]models.Status, error) {
	body := makeBody(ctx, utils.Map{
		"withFiles": onlyMedia,
		"limit":     limit,
	})
	if minId != "" {
		body["sinceId"] = minId
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	var u string
	switch timelineType {
	case models.TimelinePublicTypeLocal:
		u = "/api/notes/local-timeline"
	case models.TimelinePublicTypeRemote:
		u = "/api/notes/global-timeline"
	default:
		err := errors.New("invalid timeline type")
		return nil, err
	}
	var result []models.MkNote
	_, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post(u)
	if err != nil {
		return nil, err
	}
	list := slice.Map(result, func(_ int, n models.MkNote) models.Status { return n.ToStatus(ctx) })
	return list, nil
}

func TimelineHome(ctx misstodon.Context,
	limit int, maxId, minId string) ([]models.Status, error) {
	body := makeBody(ctx, utils.Map{"limit": limit})
	if minId != "" {
		body["sinceId"] = minId
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	var result []models.MkNote
	_, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/timeline")
	if err != nil {
		return nil, err
	}
	list := slice.Map(result, func(_ int, n models.MkNote) models.Status { return n.ToStatus(ctx) })
	return list, nil
}

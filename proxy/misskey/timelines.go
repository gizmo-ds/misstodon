package misskey

import (
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func TimelinePublic(ctx Context,
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
		u = utils.JoinURL(ctx.Server(), "/api/notes/local-timeline")
	case models.TimelinePublicTypeRemote:
		u = utils.JoinURL(ctx.Server(), "/api/notes/global-timeline")
	default:
		err := errors.New("invalid timeline type")
		return nil, err
	}
	var result []models.MkNote
	_, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(u)
	if err != nil {
		return nil, err
	}
	list := lo.Map(result, func(n models.MkNote, _ int) models.Status { return n.ToStatus(ctx.Server()) })
	return list, nil
}

func TimelineHome(ctx Context,
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
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/timeline"))
	if err != nil {
		return nil, err
	}
	list := lo.Map(result, func(n models.MkNote, _ int) models.Status { return n.ToStatus(ctx.Server()) })
	return list, nil
}

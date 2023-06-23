package misskey

import (
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func TimelinePublic(server, token string,
	timelineType models.TimelinePublicType, onlyMedia bool,
	limit int, maxId, minId string) ([]models.Status, error) {
	values := utils.Map{}
	if minId != "" {
		values["sinceId"] = minId
	}
	if maxId != "" {
		values["untilId"] = maxId
	}
	values["withFiles"] = onlyMedia
	values["limit"] = limit
	if token != "" {
		values["i"] = token
	}
	var u string
	switch timelineType {
	case models.TimelinePublicTypeLocal:
		u = utils.JoinURL(server, "/api/notes/local-timeline")
	case models.TimelinePublicTypeRemote:
		u = utils.JoinURL(server, "/api/notes/global-timeline")
	default:
		err := errors.New("invalid timeline type")
		return nil, err
	}
	var result []models.MkNote
	_, err := client.R().
		SetBody(values).
		SetResult(&result).
		Post(u)
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(server))
	}
	return list, nil
}

func TimelineHome(server, token string,
	limit int, maxId, minId string) ([]models.Status, error) {
	body := utils.Map{}
	if minId != "" {
		body["sinceId"] = minId
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	body["limit"] = limit
	if token != "" {
		body["i"] = token
	}
	var result []models.MkNote
	_, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/notes/timeline"))
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(server))
	}
	return list, nil
}

func TimelineHashtag(server, token string,
	hashtag string,
	limit int, maxId, sinceId, minId string) ([]models.Status, error) {
	body := utils.Map{"limit": limit}
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	if token != "" {
		body["i"] = token
	}
	var result []models.MkNote
	_, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/notes/search-by-tag"))
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(server))
	}
	return list, nil
}

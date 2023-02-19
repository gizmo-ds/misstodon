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
		u = "https://" + server + "/api/notes/local-timeline"
	case models.TimelinePublicTypeRemote:
		u = "https://" + server + "/api/notes/global-timeline"
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

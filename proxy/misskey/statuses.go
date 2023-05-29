package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func StatusSingle(server, token, statusID string) (models.Status, error) {
	var status models.Status
	var mkStatus models.MkNote
	body := map[string]any{
		"noteId": statusID,
	}
	if token != "" {
		body["i"] = token
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&mkStatus).
		Post(utils.JoinURL(server, "/api/notes/show"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err := isucceed(resp, 200); err != nil {
		return status, errors.WithStack(err)
	}
	status = mkStatus.ToStatus(server)
	return status, err
}

func StatusBookmark(server, token, id string) (models.Status, error) {
	status, err := StatusSingle(server, token, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(utils.Map{
			"i":      token,
			"noteId": id,
		}).
		Post(utils.JoinURL(server, "/api/notes/favorites/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusNoContent, "ALREADY_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = true
	return status, nil
}

func StatusUnBookmark(server, token, id string) (models.Status, error) {
	status, err := StatusSingle(server, token, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(utils.Map{
			"i":      token,
			"noteId": id,
		}).
		Post(utils.JoinURL(server, "/api/notes/favorites/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusNoContent, "NOT_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = false
	return status, nil
}

func StatusBookmarks(server, token string,
	limit int, sinceID, minID, maxID string) ([]models.Status, error) {
	var result []struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		Note      models.MkNote `json:"note"`
	}
	body := utils.Map{"i": token, "limit": limit}
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/i/favorites"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var status []models.Status
	for _, s := range result {
		status = append(status, s.Note.ToStatus(server))
	}
	return status, nil
}

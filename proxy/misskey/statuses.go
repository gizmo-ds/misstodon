package misskey

import (
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
		Post("https://" + server + "/api/notes/show")
	if err != nil {
		return status, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return status, ErrNotFound
	}
	status = mkStatus.ToStatus(server)
	return status, err
}

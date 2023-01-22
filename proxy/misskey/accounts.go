package misskey

import (
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAcctIsInvalid = errors.New("acct format is invalid")
)

func Lookup(server string, acct string) (models.Account, error) {
	var host *string
	var info models.Account
	username, _host := utils.AcctInfo(acct)
	if username == "" {
		return info, ErrAcctIsInvalid
	}
	if _host == "" {
		_host = server
	}
	if _host != server {
		host = &_host
	}
	var serverInfo models.MkUser
	resp, err := client.R().
		SetBody(map[string]any{
			"username": username,
			"host":     host,
		}).
		SetResult(&serverInfo).
		Post("https://" + server + "/api/users/show")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return info, ErrNotFound
	}
	return serverInfo.ToAccount(acct, server)
}

func AccountsStatuses(
	server, accountID string,
	limit int,
	pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs bool,
	maxID, minID string) ([]models.Status, error) {
	var notes []models.MkNote
	r := map[string]any{
		"userId":         accountID,
		"limit":          limit,
		"includeReplies": !excludeReplies,
	}
	if onlyMedia {
		r["fileType"] = SupportedMimeTypes
	}
	resp, err := client.R().
		SetBody(r).
		SetResult(&notes).
		Post("https://" + server + "/api/users/notes")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to get statuses")
	}
	var statuses []models.Status
	for _, note := range notes {
		statuses = append(statuses, note.ToStatus(server))
	}
	return statuses, nil
}

func VerifyCredentials(server, accessToken string) (models.Account, error) {
	var info models.Account
	var serverInfo models.MkUser
	resp, err := client.R().
		SetBody(map[string]any{
			"i": accessToken,
		}).
		SetResult(&serverInfo).
		Post("https://" + server + "/api/i")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return info, errors.New("failed to verify credentials")
	}
	return serverInfo.ToAccount(serverInfo.Username+"@"+server, server)
}

package misskey

import (
	"time"

	"github.com/gizmo-ds/misstodon/internal/mfm"
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
		return info, err
	}
	if resp.StatusCode() != 200 {
		return info, ErrNotFound
	}
	createdAt, err := time.Parse(time.RFC3339, serverInfo.CreatedAt)
	if err != nil {
		return info, err
	}
	info = models.Account{
		ID:             serverInfo.ID,
		Username:       serverInfo.Username,
		Acct:           username + "@" + _host,
		DisplayName:    serverInfo.Name,
		Locked:         serverInfo.IsLocked,
		Bot:            serverInfo.IsBot,
		CreatedAt:      createdAt.Format("2006-01-02"),
		Url:            "https://" + _host + "/@" + username,
		Avatar:         serverInfo.AvatarUrl,
		AvatarStatic:   serverInfo.AvatarUrl,
		Header:         serverInfo.BannerUrl,
		HeaderStatic:   serverInfo.BannerUrl,
		FollowersCount: serverInfo.FollowersCount,
		FollowingCount: serverInfo.FollowingCount,
		StatusesCount:  serverInfo.NotesCount,
		Emojis:         []models.CustomEmoji{},
		Fields:         serverInfo.Fields,
	}
	_lastStatusAt := serverInfo.UpdatedAt
	if serverInfo.UpdatedAt != nil {
		lastStatusAt, err := time.Parse(time.RFC3339, *_lastStatusAt)
		if err != nil {
			return info, err
		}
		t := lastStatusAt.Format("2006-01-02")
		info.LastStatusAt = &t
	}
	if serverInfo.Description != nil {
		info.Note, err = mfm.ToHtml(*serverInfo.Description, mfm.Option{Url: "https://" + server})
		if err != nil {
			return info, err
		}
	}
	return info, nil
}

func AccountsStatuses(
	server, accountID string,
	limit int,
	pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs bool,
	maxID, minID string) {

}

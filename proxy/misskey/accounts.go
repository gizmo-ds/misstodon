package misskey

import (
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
)

func Lookup(server string, acct string) (models.Account, error) {
	var host *string
	username, _host := utils.AcctInfo(acct)
	if _host == "" {
		_host = server
	}
	if _host != server {
		host = &_host
	}
	var info models.Account
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
		return info, err
	}
	createdAt, err := time.Parse(time.RFC3339, serverInfo.CreatedAt)
	if err != nil {
		return info, err
	}
	_lastStatusAt := serverInfo.UpdatedAt
	if _lastStatusAt != nil {
		lastStatusAt, err := time.Parse(time.RFC3339, *_lastStatusAt)
		if err != nil {
			return info, err
		}
		t := lastStatusAt.Format("2006-01-02")
		_lastStatusAt = &t
	}
	return models.Account{
		ID:             serverInfo.ID,
		Username:       serverInfo.Username,
		Acct:           username + "@" + _host,
		DisplayName:    serverInfo.Name,
		Locked:         serverInfo.IsLocked,
		Bot:            serverInfo.IsBot,
		CreatedAt:      createdAt.Format("2006-01-02"),
		LastStatusAt:   _lastStatusAt,
		Note:           utils.MfmToHtml(serverInfo.Description),
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
	}, nil
}

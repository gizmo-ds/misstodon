package misskey

import (
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func OAuthAuthorize(server, secret string) (string, error) {
	var result struct {
		Token string `json:"token"`
		Url   string `json:"url"`
	}
	resp, err := client.R().
		SetBody(map[string]any{
			"appSecret": secret,
		}).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/auth/session/generate"))
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return "", errors.New("failed to authorize")
	}
	return result.Url, nil
}

func OAuthToken(server, token, secret string) (string, error) {
	var result struct {
		AccessToken string `json:"accessToken"`
		User        models.MkUser
	}
	resp, err := client.R().
		SetBody(map[string]any{
			"appSecret": secret,
			"token":     token,
		}).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/auth/session/userkey"))
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return "", errors.New("failed to get access_token")
	}
	return result.AccessToken, nil
}

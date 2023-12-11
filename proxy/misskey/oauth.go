package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func OAuthAuthorize(ctx misstodon.Context, secret string) (string, error) {
	var result struct {
		Token string `json:"token"`
		Url   string `json:"url"`
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(map[string]any{
			"appSecret": secret,
		}).
		SetResult(&result).
		Post("/api/auth/session/generate")
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("failed to authorize")
	}
	return result.Url, nil
}

func OAuthToken(ctx misstodon.Context, token, secret string) (string, string, error) {
	var result struct {
		AccessToken string `json:"accessToken"`
		User        models.MkUser
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(map[string]any{
			"appSecret": secret,
			"token":     token,
		}).
		SetResult(&result).
		Post("/api/auth/session/userkey")
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", "", errors.New("failed to get access_token")
	}
	return result.AccessToken, result.User.ID, nil
}

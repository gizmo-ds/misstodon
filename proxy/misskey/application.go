package misskey

import (
	"net/http"
	"strings"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func ApplicationCreate(ctx misstodon.Context, clientName, redirectUris, scopes, website string) (models.Application, error) {
	var permissions []string
	var app models.Application
	arr := strings.Split(scopes, " ")
	for _, scope := range arr {
		switch scope {
		case "read":
			permissions = append(permissions, models.ApplicationPermissionRead...)
		case "write":
			permissions = append(permissions, models.ApplicationPermissionWrite...)
		case "follow":
			permissions = append(permissions, models.ApplicationPermissionFollow...)
		case "push":
			// FIXME: 未实现WebPushAPI
		default:
			permissions = append(permissions, scope)
		}
	}
	permissions = utils.Unique(permissions)
	var result models.MkApplication
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(map[string]any{
			"name":        clientName,
			"description": website,
			"callbackUrl": redirectUris,
			"permission":  permissions,
		}).
		SetResult(&result).
		Post("/api/app/create")
	if err != nil {
		return app, err
	}
	if resp.StatusCode() != http.StatusOK {
		return app, errors.New("failed to create application")
	}
	app = models.Application{
		ID:           result.ID,
		Name:         result.Name,
		RedirectUri:  result.CallbackUrl,
		ClientID:     &result.ID,
		ClientSecret: &result.Secret,
		// FIXME: 未实现WebPushAPI
		VapidKey: "",
	}
	return app, nil
}

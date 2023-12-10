package v1

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func ApplicationRouter(e *echo.Group) {
	group := e.Group("/apps")
	group.POST("", ApplicationCreateHandler)
}

func ApplicationCreateHandler(c echo.Context) error {
	var params struct {
		ClientName   string `json:"client_name" form:"client_name"`
		WebSite      string `json:"website" form:"website"`
		RedirectUris string `json:"redirect_uris" form:"redirect_uris"`
		Scopes       string `json:"scopes" form:"scopes"`
	}
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if params.ClientName == "" || params.RedirectUris == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "client_name and redirect_uris are required")
	}

	ctx, err := misstodon.ContextWithEchoContext(c, false)
	if err != nil {
		return err
	}

	u, err := url.Parse(strings.Join([]string{"https://", c.Request().Host, "/oauth/redirect"}, ""))
	if err != nil {
		return errors.WithStack(err)
	}
	query := u.Query()
	query.Add("server", ctx.ProxyServer())
	query.Add("redirect_uris", params.RedirectUris)
	u.RawQuery = query.Encode()

	app, err := misskey.ApplicationCreate(
		ctx,
		params.ClientName,
		u.String(),
		params.Scopes,
		params.WebSite)
	if err != nil {
		return err
	}
	err = global.DB.Set(ctx.ProxyServer(), app.ID, *app.ClientSecret, -1)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, app)
}

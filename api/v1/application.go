package v1

import (
	"net/http"
	"net/url"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func ApplicationRouter(e *echo.Group) {
	group := e.Group("/apps")
	group.POST("", ApplicationCreate)
}

func ApplicationCreate(c echo.Context) error {
	var params struct {
		ClientName   string `json:"client_name"`
		WebSite      string `json:"website"`
		RedirectUris string `json:"redirect_uris"`
		Scopes       string `json:"scopes"`
	}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	if params.ClientName == "" || params.RedirectUris == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "client_name and redirect_uris are required",
		})
	}
	server := c.Get("server").(string)
	u, err := url.Parse(global.Config.Server.Url + "/oauth/redirect")
	if err != nil {
		return errors.WithStack(err)
	}
	query := u.Query()
	query.Add("server", server)
	query.Add("redirect_uris", params.RedirectUris)
	u.RawQuery = query.Encode()
	app, err := misskey.ApplicationCreate(
		server,
		params.ClientName,
		u.String(),
		params.Scopes,
		params.WebSite)
	if err != nil {
		return err
	}
	err = global.DB.Set(server, app.ID, *app.ClientSecret, -1)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, app)
}

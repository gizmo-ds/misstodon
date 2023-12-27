package oauth

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Group) {
	group := e.Group("/oauth", middleware.CORS())
	group.GET("/authorize", AuthorizeHandler)
	group.POST("/token", TokenHandler)
	// NOTE: This is not a standard endpoint
	group.GET("/redirect", RedirectHandler)
}

func RedirectHandler(c echo.Context) error {
	redirectUris := c.QueryParam("redirect_uris")
	server := c.QueryParam("server")
	token := c.QueryParam("token")
	if redirectUris == "" || server == "" {
		return c.String(http.StatusBadRequest, "redirect_uris and server are required")
	}
	if token == "" {
		if strings.Contains(redirectUris, "?token=") {
			i := strings.Index(redirectUris, "?token=")
			token = redirectUris[i+7:]
			redirectUris = redirectUris[:i]
		}
		if strings.Contains(server, "?token=") {
			i := strings.Index(server, "?token=")
			token = server[i+7:]
			server = server[:i]
		}
	}
	u, err := url.Parse(redirectUris)
	if err != nil {
		return err
	}
	query := u.Query()
	query.Add("code", token)
	u.RawQuery = query.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

func TokenHandler(c echo.Context) error {
	var params struct {
		GrantType    string `json:"grant_type" form:"grant_type"`
		ClientID     string `json:"client_id" form:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret"`
		RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
		Code         string `json:"code" form:"code"`
		Scope        string `json:"scope" form:"scope"`
	}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	if params.GrantType == "" || params.ClientID == "" ||
		params.ClientSecret == "" || params.RedirectURI == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "grant_type, client_id, client_secret and redirect_uri are required",
		})
	}
	ctx, err := misstodon.ContextWithEchoContext(c, false)
	if err != nil {
		return err
	}
	accessToken, userID, err := misskey.OAuthToken(ctx, params.Code, params.ClientSecret)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token": strings.Join([]string{userID, accessToken}, "."),
		"token_type":   "Bearer",
		"scope":        params.Scope,
		"created_at":   time.Now().Unix(),
	})
}

func AuthorizeHandler(c echo.Context) error {
	var params struct {
		ClientID     string `query:"client_id"`
		RedirectUri  string `query:"redirect_uri"`
		ResponseType string `query:"response_type"`
		Scope        string `query:"scope"`
		Lang         string `query:"lang"`
		ForceLogin   bool   `query:"force_login"`
	}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	if params.ResponseType != "code" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "response_type must be code",
		})
	}
	if params.ClientID == "" || params.RedirectUri == "" || params.ResponseType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "client_id, redirect_uri and response_type are required",
		})
	}
	ctx, err := misstodon.ContextWithEchoContext(c, false)
	if err != nil {
		return err
	}
	secret, ok := global.DB.Get(ctx.ProxyServer(), params.ClientID)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "client_id is invalid",
		})
	}
	u, err := misskey.OAuthAuthorize(ctx, secret)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, u)
}

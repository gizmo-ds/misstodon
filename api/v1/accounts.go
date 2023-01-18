package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func AccountsVerifyCredentials(c echo.Context) error {
	return nil
}

func AccountsLookup(c echo.Context) error {
	acct := c.QueryParam("acct")
	if acct == "" {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: "acct is required",
		})
	}
	info, err := misskey.Lookup(c.Get("server").(string), acct)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
		if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{
				Error: "Record not found",
			})
		} else if errors.Is(err, misskey.ErrAcctIsInvalid) {
			return c.JSON(http.StatusBadRequest, httperror.ServerError{
				Error: err.Error(),
			})
		}
		return err
	}
	return c.JSON(http.StatusOK, info)
}

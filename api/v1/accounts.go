package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

func AccountsRouter(e *echo.Group) {
	group := e.Group("/accounts", middleware.CORS())
	group.GET("/verify_credentials", AccountsVerifyCredentialsHandler)
	group.GET("/lookup", AccountsLookupHandler)
	group.GET("/:accountID/statuses", AccountsStatusesHandler)
}

func AccountsVerifyCredentialsHandler(c echo.Context) error {
	accessToken, err := utils.GetHeaderToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: err.Error(),
		})
	}
	server := c.Get("server").(string)
	info, err := misskey.VerifyCredentials(server, accessToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func AccountsLookupHandler(c echo.Context) error {
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

func AccountsStatusesHandler(c echo.Context) error {
	accountID := c.Param("accountID")
	limit := 30
	pinnedOnly := false
	onlyMedia := false
	onlyPublic := false
	excludeReplies := false
	excludeReblogs := false
	maxID := ""
	minID := ""
	if err := echo.QueryParamsBinder(c).
		Int("limit", &limit).
		Bool("pinned_only", &pinnedOnly).
		Bool("only_media", &onlyMedia).
		Bool("only_public", &onlyPublic).
		Bool("exclude_replies", &excludeReplies).
		Bool("exclude_reblogs", &excludeReblogs).
		String("max_id", &maxID).
		String("min_id", &minID).
		BindError(); err != nil {
		e := err.(*echo.BindingError)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"field": e.Field,
			"error": e.Message,
		})
	}
	statuses, err := misskey.AccountsStatuses(
		c.Get("server").(string), accountID,
		limit,
		pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs,
		maxID, minID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, statuses)
}

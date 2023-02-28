package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func AccountsRouter(e *echo.Group) {
	group := e.Group("/accounts")
	group.GET("/verify_credentials", AccountsVerifyCredentialsHandler)
	group.PATCH("/update_credentials", AccountsUpdateCredentialsHandler)
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
	if info.Header == "" || info.HeaderStatic == "" {
		info.Header = fmt.Sprintf("%s://%s/static/missing.png", c.Scheme(), c.Request().Host)
		info.HeaderStatic = info.Header
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
	return c.JSON(http.StatusOK, utils.SliceIfNull(statuses))
}

func AccountsUpdateCredentialsHandler(c echo.Context) error {
	form, err := parseAccountsUpdateCredentialsForm(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}
	server := c.Get("server").(string)
	accessToken, err := utils.GetHeaderToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: err.Error(),
		})
	}
	account, err := misskey.UpdateCredentials(server, accessToken,
		form.DisplayName, form.Note,
		form.Locked, form.Bot, form.Discoverable,
		form.SourcePrivacy, form.SourceSensitive, form.SourceLanguage,
		form.AccountFields,
		form.Avatar, form.Header)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, account)
}

type accountsUpdateCredentialsForm struct {
	DisplayName     *string `form:"display_name"`
	Note            *string `form:"note"`
	Locked          *bool   `form:"locked"`
	Bot             *bool   `form:"bot"`
	Discoverable    *bool   `form:"discoverable"`
	SourcePrivacy   *string `form:"source[privacy]"`
	SourceSensitive *bool   `form:"source[sensitive]"`
	SourceLanguage  *string `form:"source[language]"`
	AccountFields   []models.AccountField
	Avatar          *multipart.FileHeader
	Header          *multipart.FileHeader
}

func parseAccountsUpdateCredentialsForm(c echo.Context) (f accountsUpdateCredentialsForm, err error) {
	var form accountsUpdateCredentialsForm
	if err = c.Bind(&form); err != nil {
		return
	}

	var values = make(map[string][]string)
	for k, v := range c.QueryParams() {
		values[k] = v
	}
	if fp, err := c.FormParams(); err == nil {
		for k, v := range fp {
			values[k] = v
		}
	}
	if mf, err := c.MultipartForm(); err == nil {
		for k, v := range mf.Value {
			values[k] = v
		}
	}
	for _, field := range utils.GetFieldsAttributes(values) {
		form.AccountFields = append(form.AccountFields, models.AccountField(field))
	}
	if fh, err := c.FormFile("avatar"); err == nil {
		form.Avatar = fh
	}
	if fh, err := c.FormFile("header"); err == nil {
		form.Header = fh
	}
	return form, nil
}

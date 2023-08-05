package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func AccountsRouter(e *echo.Group) {
	group := e.Group("/accounts")
	e.GET("/favourites", AccountFavourites)
	group.GET("/verify_credentials", AccountsVerifyCredentialsHandler)
	group.PATCH("/update_credentials", AccountsUpdateCredentialsHandler)
	group.GET("/lookup", AccountsLookupHandler)
	group.GET("/:id", AccountsGetHandler)
	group.GET("/:id/statuses", AccountsStatusesHandler)
	group.GET("/:id/followers", AccountFollowers)
	group.GET("/:id/following", AccountFollowing)
	group.GET("/relationships", AccountRelationships)
	group.POST("/:id/follow", AccountFollow)
	group.POST("/:id/unfollow", AccountUnfollow)
	group.POST("/:id/mute", AccountMute)
	group.POST("/:id/unmute", AccountUnmute)
}

func AccountsVerifyCredentialsHandler(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	info, err := misskey.VerifyCredentials(ctx)
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
	ctx, _ := misstodon.ContextWithEchoContext(c)
	info, err := misskey.AccountsLookup(ctx, acct)
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
	uid := c.Param("id")

	ctx, _ := misstodon.ContextWithEchoContext(c)

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
		var e *echo.BindingError
		errors.As(err, &e)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"field": e.Field,
			"error": e.Message,
		})
	}
	statuses, err := misskey.AccountsStatuses(
		ctx, uid,
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

	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}

	account, err := misskey.UpdateCredentials(ctx,
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

func AccountFollowRequests(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	var query struct {
		Limit   int    `query:"limit"`
		MaxID   string `query:"max_id"`
		SinceID string `query:"since_id"`
	}
	if err = c.Bind(&query); err != nil {
		return err
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	accounts, err := misskey.AccountFollowRequests(ctx, query.Limit, query.SinceID, query.MaxID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountFollowers(c echo.Context) error {
	ctx, _ := misstodon.ContextWithEchoContext(c)
	id := c.Param("id")
	var query struct {
		Limit   int    `query:"limit"`
		MaxID   string `query:"max_id"`
		MinID   string `query:"min_id"`
		SinceID string `query:"since_id"`
	}
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	if query.Limit > 80 {
		query.Limit = 80
	}
	accounts, err := misskey.AccountFollowers(ctx, id, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountFollowing(c echo.Context) error {
	ctx, _ := misstodon.ContextWithEchoContext(c)

	id := c.Param("id")
	var query struct {
		Limit   int    `query:"limit"`
		MaxID   string `query:"max_id"`
		MinID   string `query:"min_id"`
		SinceID string `query:"since_id"`
	}
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	if query.Limit > 80 {
		query.Limit = 80
	}
	accounts, err := misskey.AccountFollowing(ctx, id, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountRelationships(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	var ids []string
	for k, v := range c.QueryParams() {
		if k == "id[]" {
			ids = append(ids, v...)
			continue
		}
	}
	relationships, err := misskey.AccountRelationships(ctx, ids)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relationships)
}

func AccountFollow(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	id := c.Param("id")
	if err = misskey.AccountFollow(ctx, id); err != nil {
		return err
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relationships[0])
}

func AccountUnfollow(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	id := c.Param("id")
	if err = misskey.AccountUnfollow(ctx, id); err != nil {
		return err
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relationships[0])
}

func AccountMute(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	var params struct {
		ID       string `param:"id"`
		Duration int64  `json:"duration" form:"duration"`
	}
	if err := c.Bind(&params); err != nil {
		return err
	}
	if err = misskey.AccountMute(ctx, params.ID, params.Duration); err != nil {
		return err
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{params.ID})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relationships[0])
}

func AccountUnmute(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	id := c.Param("id")
	if err = misskey.AccountUnmute(ctx, id); err != nil {
		return err
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relationships[0])
}

func AccountsGetHandler(c echo.Context) error {
	ctx, _ := misstodon.ContextWithEchoContext(c)
	info, err := misskey.AccountsGet(ctx, c.Param("id"))
	if err != nil {
		return err
	}
	if info.Header == "" || info.HeaderStatic == "" {
		info.Header = fmt.Sprintf("%s://%s/static/missing.png", c.Scheme(), c.Request().Host)
		info.HeaderStatic = info.Header
	}
	return c.JSON(http.StatusOK, info)
}

func AccountFavourites(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}

	var params struct {
		Limit   int    `query:"limit"`
		MaxID   string `query:"max_id"`
		MinID   string `query:"min_id"`
		SinceID string `query:"since_id"`
	}
	if err = c.Bind(&params); err != nil {
		return err
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	list, err := misskey.AccountFavourites(ctx,
		params.Limit, params.SinceID, params.MinID, params.MaxID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}

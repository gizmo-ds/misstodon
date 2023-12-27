package v1

import (
	"net/http"
	"time"

	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func StatusesRouter(e *echo.Group) {
	group := e.Group("/statuses")
	group.POST("", PostNewStatusHandler)
	group.GET("/:id", StatusHandler)
	group.DELETE("/:id", StatusDeleteHandler)
	group.GET("/:id/context", StatusContextHandler)
	group.POST("/:id/reblog", StatusReblogHandler)
	group.POST("/:id/bookmark", StatusBookmarkHandler)
	group.POST("/:id/unbookmark", StatusUnBookmarkHandler)
	group.POST("/:id/favourite", StatusFavouriteHandler)
	group.POST("/:id/unfavourite", StatusUnFavouriteHandler)
}

func StatusHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithEchoContext(c)
	info, err := misskey.StatusSingle(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func StatusContextHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, false)
	if err != nil {
		return err
	}
	result, err := misskey.StatusContext(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func StatusFavouriteHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusFavourite(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

func StatusUnFavouriteHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusUnFavourite(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

func StatusBookmarkHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusBookmark(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

func StatusUnBookmarkHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusUnBookmark(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else if errors.Is(err, misskey.ErrNotFound) {
			return c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, status)
}

func StatusBookmarks(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	var query struct {
		Limit   int    `query:"limit"`
		MaxID   string `query:"max_id"`
		MinID   string `query:"min_id"`
		SinceID string `query:"since_id"`
	}
	if err = c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	status, err := misskey.StatusBookmarks(ctx, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
		} else {
			return err
		}
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(status))
}

type postNewStatusForm struct {
	Status      *string                 `json:"status"`
	Poll        any                     `json:"poll"` // FIXME: Poll 未实现
	MediaIDs    []string                `json:"media_ids"`
	InReplyToID string                  `json:"in_reply_to_id"`
	Sensitive   bool                    `json:"sensitive"`
	SpoilerText string                  `json:"spoiler_text"`
	Visibility  models.StatusVisibility `json:"visibility"`
	Language    string                  `json:"language"`
	ScheduledAt time.Time               `json:"scheduled_at"`
}

func PostNewStatusHandler(c echo.Context) error {
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}

	var form postNewStatusForm
	if err = c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
	}
	status, err := misskey.PostNewStatus(ctx,
		form.Status, form.Poll, form.MediaIDs, form.InReplyToID,
		form.Sensitive, form.SpoilerText,
		form.Visibility, form.Language,
		form.ScheduledAt)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, status)
}

func StatusReblogHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusReblog(ctx, id, models.StatusVisibilityPublic)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, status)
}

func StatusDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithEchoContext(c, true)
	if err != nil {
		return err
	}
	status, err := misskey.StatusDelete(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, status)
}

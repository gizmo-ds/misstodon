package misskey

import (
	"net/http"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func StatusSingle(ctx Context, statusID string) (models.Status, error) {
	var status models.Status
	var note models.MkNote
	body := makeBody(ctx, utils.Map{"noteId": statusID})
	resp, err := client.R().
		SetBody(body).
		SetResult(&note).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/show"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, 200); err != nil {
		return status, errors.WithStack(err)
	}
	status = note.ToStatus(ctx.Server())
	if ctx.Token() != nil {
		state, err := getNoteState(ctx.Server(), *ctx.Token(), status.ID)
		if err != nil {
			return status, err
		}
		status.Bookmarked = state.IsFavorited
		status.Muted = state.IsMutedThread
	}
	return status, err
}

type noteState struct {
	IsFavorited   bool `json:"isFavorited"`
	IsMutedThread bool `json:"isMutedThread"`
}

func getNoteState(server, token, noteId string) (noteState, error) {
	var state noteState
	resp, err := client.R().
		SetBody(utils.Map{"i": token, "noteId": noteId}).
		SetResult(&state).
		Post(utils.JoinURL(server, "/api/notes/state"))
	if err != nil {
		return state, errors.WithStack(err)
	}
	if err = isucceed(resp, 200); err != nil {
		return state, errors.WithStack(err)
	}
	return state, nil
}

func StatusFavourite(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{
			"noteId":   id,
			"reaction": "⭐",
		})).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/reactions/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_REACTED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Favourited = true
	status.FavouritesCount += 1
	return status, nil
}

func StatusUnFavourite(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/reactions/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "NOT_REACTED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Favourited = false
	status.FavouritesCount -= 1
	return status, nil
}

func StatusBookmark(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/favorites/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = true
	return status, nil
}

func StatusUnBookmark(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/favorites/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "NOT_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = false
	return status, nil
}

// StatusBookmarks
// NOTE: 为了减少请求数量, 不支持 Bookmarked
func StatusBookmarks(ctx Context,
	limit int, sinceID, minID, maxID string) ([]models.Status, error) {
	var result []struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		Note      models.MkNote `json:"note"`
	}
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.Server(), "/api/i/favorites"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var status []models.Status
	for _, s := range result {
		status = append(status, s.Note.ToStatus(ctx.Server()))
	}
	return status, nil
}

// PostNewStatus 发送新的 Status
// FIXME: Poll 未实现
func PostNewStatus(ctx Context,
	status *string, Poll any, MediaIDs []string, InReplyToID string,
	Sensitive bool, SpoilerText string,
	Visibility models.StatusVisibility, Language string,
	ScheduledAt time.Time,
) (any, error) {
	body := makeBody(ctx, utils.Map{"localOnly": false})
	var noteMentions []string
	if status != nil && *status != "" {
		body["text"] = *status
		noteMentions = append(noteMentions, utils.GetMentions(*status)...)
	}
	if Sensitive {
		if SpoilerText != "" {
			body["cw"] = SpoilerText
		} else {
			body["cw"] = "Sensitive"
		}
	}
	switch Visibility {
	case models.StatusVisibilityPublic:
		body["visibility"] = "public"
	case models.StatusVisibilityUnlisted:
		body["visibility"] = "home"
	case models.StatusVisibilityPrivate:
		body["visibility"] = "followers"
	case models.StatusVisibilityDirect:
		body["visibility"] = "specified"
		var visibleUserIds []string
		for _, m := range noteMentions {
			a, err := AccountsLookup(ctx, m)
			if err != nil {
				return nil, err
			}
			visibleUserIds = append(visibleUserIds, a.ID)
		}
		if len(visibleUserIds) > 0 {
			body["visibleUserIds"] = visibleUserIds
		}
	}
	if MediaIDs != nil {
		body["mediaIds"] = MediaIDs
	}
	if InReplyToID != "" {
		body["replyId"] = InReplyToID
	}
	var result struct {
		CreatedNote models.MkNote `json:"createdNote"`
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/create"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return result.CreatedNote.ToStatus(ctx.Server()), nil
}

func SearchStatusByHashtag(ctx Context,
	hashtag string,
	limit int, maxId, sinceId, minId string) ([]models.Status, error) {
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	var result []models.MkNote
	_, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.Server(), "/api/notes/search-by-tag"))
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(ctx.Server()))
	}
	return list, nil
}

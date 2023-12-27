package misskey

import (
	"net/http"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/pkg/errors"
)

func StatusSingle(ctx misstodon.Context, statusID string) (models.Status, error) {
	var status models.Status
	var note models.MkNote
	body := makeBody(ctx, utils.Map{"noteId": statusID})
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&note).
		Post("/api/notes/show")
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return status, errors.WithStack(err)
	}
	status = note.ToStatus(ctx)
	if ctx.Token() != nil {
		state, err := getNoteState(ctx.ProxyServer(), *ctx.Token(), status.ID)
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
		SetBaseURL(server).
		SetBody(utils.Map{"i": token, "noteId": noteId}).
		SetResult(&state).
		Post("/api/notes/state")
	if err != nil {
		return state, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return state, errors.WithStack(err)
	}
	return state, nil
}

func StatusFavourite(ctx misstodon.Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, utils.Map{
			"noteId":   id,
			"reaction": "⭐",
		})).
		Post("/api/notes/reactions/create")
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

func StatusUnFavourite(ctx misstodon.Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post("/api/notes/reactions/delete")
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

func StatusBookmark(ctx misstodon.Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post("/api/notes/favorites/create")
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = true
	return status, nil
}

func StatusUnBookmark(ctx misstodon.Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post("/api/notes/favorites/delete")
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
func StatusBookmarks(ctx misstodon.Context,
	limit int, sinceId, minId, maxId string) ([]models.Status, error) {
	type favorite struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		Note      models.MkNote `json:"note"`
	}
	var result []favorite
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/i/favorites")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	status := slice.Map(result, func(_ int, item favorite) models.Status {
		s := item.Note.ToStatus(ctx)
		s.Bookmarked = true
		return s
	})
	return status, nil
}

// PostNewStatus 发送新的 Status
// FIXME: Poll 未实现
func PostNewStatus(ctx misstodon.Context,
	status *string, poll any, mediaIds []string, inReplyToId string,
	sensitive bool, spoilerText string,
	visibility models.StatusVisibility, language string,
	scheduledAt time.Time,
) (models.Status, error) {
	body := makeBody(ctx, utils.Map{"localOnly": false})
	var noteMentions []string
	if status != nil && *status != "" {
		body["text"] = *status
		noteMentions = append(noteMentions, utils.GetMentions(*status)...)
	}
	if sensitive {
		if spoilerText != "" {
			body["cw"] = spoilerText
		} else {
			body["cw"] = "Sensitive"
		}
	}
	body["visibility"] = visibility.ToMkNoteVisibility()
	if visibility == models.StatusVisibilityDirect {
		var visibleUserIds []string
		for _, m := range noteMentions {
			a, err := AccountsLookup(ctx, m)
			if err != nil {
				return models.Status{}, err
			}
			visibleUserIds = append(visibleUserIds, a.ID)
		}
		if len(visibleUserIds) > 0 {
			body["visibleUserIds"] = visibleUserIds
		}
	}
	if len(mediaIds) > 0 {
		body["mediaIds"] = mediaIds
	}
	if inReplyToId != "" {
		body["replyId"] = inReplyToId
	}
	var result struct {
		CreatedNote models.MkNote `json:"createdNote"`
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/create")
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	return result.CreatedNote.ToStatus(ctx), nil
}

func SearchStatusByHashtag(ctx misstodon.Context,
	hashtag string,
	limit int, maxId, sinceId, minId string) ([]models.Status, error) {
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	body["tag"] = hashtag
	var result []models.MkNote
	_, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/search-by-tag")
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(ctx))
	}
	return list, nil
}

func StatusContext(ctx misstodon.Context, id string) (models.Context, error) {
	result := models.Context{
		Ancestors:   make([]models.Status, 0),
		Descendants: make([]models.Status, 0),
	}
	s, err := StatusSingle(ctx, id)
	if err != nil {
		return result, err
	}
	result, err = statusContext(ctx, s)
	return result, err
}

func statusContext(ctx misstodon.Context, status models.Status) (models.Context, error) {
	result := models.Context{
		Ancestors:   make([]models.Status, 0),
		Descendants: make([]models.Status, 0),
	}
	if status.RepliesCount > 0 {
		notes, err := noteChildren(ctx, status.ID, 30, "", "")
		if err == nil {
			result.Descendants = slice.Map(notes, func(_ int, item models.MkNote) models.Status {
				return item.ToStatus(ctx)
			})
		}
		lc := status.RepliesCount / 100
		if status.RepliesCount%100 > 0 {
			lc++
		}
		for i := 0; i < lc; i++ {
			notes, err = noteConversation(ctx, status.ID, 100, i*100)
			if err == nil {
				arr := slice.Map(notes, func(_ int, item models.MkNote) models.Status {
					return item.ToStatus(ctx)
				})
				result.Ancestors = append(result.Ancestors, arr...)
			}
		}
	}
	return result, nil
}

func noteConversation(ctx misstodon.Context, id string, limit, offset int) ([]models.MkNote, error) {
	body := makeBody(ctx, utils.Map{
		"limit":  limit,
		"offset": offset,
		"noteId": id,
	})
	var result []models.MkNote
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/conversation")
	if err != nil {
		return nil, err
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func noteChildren(ctx misstodon.Context, id string, limit int, sinceId, untilId string) ([]models.MkNote, error) {
	body := makeBody(ctx, utils.Map{
		"limit":  limit,
		"noteId": id,
	})
	if sinceId != "" {
		body["sinceId"] = sinceId
	}
	if untilId != "" {
		body["untilId"] = untilId
	}
	var result []models.MkNote
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/children")
	if err != nil {
		return nil, err
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func StatusReblog(ctx misstodon.Context, reNoteId string, visibility models.StatusVisibility) (models.Status, error) {
	body := makeBody(ctx, utils.Map{"renoteId": reNoteId})
	body["visibility"] = visibility.ToMkNoteVisibility()
	var result struct {
		CreatedNote models.MkNote `json:"createdNote"`
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/notes/create")
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	return result.CreatedNote.ToStatus(ctx), nil
}

func StatusDelete(ctx misstodon.Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, err
	}
	body := makeBody(ctx, utils.Map{"noteId": id})
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		Post("/api/notes/delete")
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent); err != nil {
		return status, errors.WithStack(err)
	}
	return status, nil
}

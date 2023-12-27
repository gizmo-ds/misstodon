package models

import (
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/mfm"
)

type MkNoteVisibility string

const (
	MkNoteVisibilityPublic MkNoteVisibility = "public"
	MkNoteVisibilityHome   MkNoteVisibility = "home"
	// MkNoteVisibilityFollow MkNoteVisibility = "follow"
	MkNoteVisibilityFollow MkNoteVisibility = "followers"
	MkNoteVisibilitySpecif MkNoteVisibility = "specified"
)

type MkNote struct {
	ID           string           `json:"id"`
	CreatedAt    string           `json:"createdAt"`
	ReplyID      *string          `json:"replyId"`
	ThreadId     *string          `json:"threadId"`
	Text         *string          `json:"text"`
	Name         *string          `json:"name"`
	Cw           *string          `json:"cw"`
	UserId       string           `json:"userId"`
	User         *MkUser          `json:"user"`
	LocalOnly    bool             `json:"localOnly"`
	Reply        *MkNote          `json:"reply"`
	ReNote       *MkNote          `json:"renote"`
	ReNoteId     *string          `json:"renoteId"`
	ReNoteCount  int              `json:"renoteCount"`
	RepliesCount int              `json:"repliesCount"`
	Reactions    map[string]int   `json:"reactions"`
	Visibility   MkNoteVisibility `json:"visibility"`
	Uri          *string          `json:"uri"`
	Url          *string          `json:"url"`
	Score        int              `json:"score"`
	FileIds      []string         `json:"fileIds"`
	Files        []MkFile         `json:"files"`
	Tags         []string         `json:"tags"`
	MyReaction   string           `json:"myReaction"`
}

func (n *MkNote) ToStatus(ctx misstodon.Context) Status {
	uid := ctx.UserID()
	s := Status{
		ID:               n.ID,
		Url:              utils.JoinURL(ctx.ProxyServer(), "/notes/", n.ID),
		Uri:              utils.JoinURL(ctx.ProxyServer(), "/notes/", n.ID),
		CreatedAt:        n.CreatedAt,
		Emojis:           []struct{}{},
		MediaAttachments: []MediaAttachment{},
		Mentions:         []StatusMention{},
		ReBlogsCount:     n.ReNoteCount,
		RepliesCount:     n.RepliesCount,
		Favourited:       n.MyReaction != "",
	}
	s.FavouritesCount = func() int {
		var count int
		for _, r := range n.Reactions {
			count += r
		}
		return count
	}()
	for _, tag := range n.Tags {
		s.Tags = append(s.Tags, StatusTag{
			Name: tag,
			Url:  utils.JoinURL(ctx.ProxyServer(), "/tags/", tag),
		})
	}
	if n.Text != nil {
		s.Content = *n.Text
		if content, err := mfm.ToHtml(*n.Text, mfm.Option{
			Url:            utils.JoinURL(ctx.ProxyServer()),
			HashtagHandler: mfm.MastodonHashtagHandler,
		}); err == nil {
			s.Content = content
		}
	}
	if n.User != nil {
		a, err := n.User.ToAccount(ctx)
		if err == nil {
			s.Account = a
		}
	}
	s.Visibility = n.Visibility.ToStatusVisibility()
	for _, file := range n.Files {
		if file.IsSensitive {
			s.Sensitive = true
		}
		// NOTE: Misskey does not return width and height of media files.
		if file.Properties.Width <= 0 && file.Properties.Height <= 0 {
			file.Properties.Width = 1920
			file.Properties.Height = 1080
		}
		s.MediaAttachments = append(s.MediaAttachments, file.ToMediaAttachment())
	}
	if n.Cw != nil {
		s.SpoilerText = *n.Cw
	}
	if n.ReNote != nil {
		re := n.ReNote.ToStatus(ctx)
		s.ReBlog = &re
	}
	if uid != nil {
		s.ReBlogged = n.ReNote != nil && n.UserId == *uid
	}
	return s
}

func (v MkNoteVisibility) ToStatusVisibility() StatusVisibility {
	switch v {
	case MkNoteVisibilityPublic:
		return StatusVisibilityPublic
	case MkNoteVisibilityHome:
		return StatusVisibilityUnlisted
	case MkNoteVisibilityFollow:
		return StatusVisibilityPrivate
	default:
		return StatusVisibilityDirect
	}
}

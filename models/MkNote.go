package models

import (
	"github.com/gizmo-ds/misstodon/internal/mfm"
	"github.com/gizmo-ds/misstodon/internal/utils"
)

type MkNoteVisibility = string

const (
	MkNoteVisibilityPublic MkNoteVisibility = "public"
	MkNoteVisibilityHome   MkNoteVisibility = "home"
	MkNoteVisibilityFollow MkNoteVisibility = "follow"
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
}

func (n *MkNote) ToStatus(server string) Status {
	s := Status{
		ID:               n.ID,
		Url:              utils.JoinURL(server, "/notes/", n.ID),
		Uri:              utils.JoinURL(server, "/notes/", n.ID),
		CreatedAt:        n.CreatedAt,
		Emojis:           []struct{}{},
		MediaAttachments: []MediaAttachment{},
		Mentions:         []StatusMention{},
		ReBlogsCount:     n.ReNoteCount,
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
			Url:  utils.JoinURL(server, "/tags/", tag),
		})
	}
	if n.Text != nil {
		s.Content = *n.Text
		if content, err := mfm.ToHtml(*n.Text, mfm.Option{
			Url:            utils.JoinURL(server),
			HashtagHandler: mfm.MastodonHashtagHandler,
		}); err == nil {
			s.Content = content
		}
		utils.GetMentions(*n.Text)
	}
	if n.User != nil {
		a, err := n.User.ToAccount(server)
		if err == nil {
			s.Account = a
		}
	}
	switch n.Visibility {
	case MkNoteVisibilityPublic, MkNoteVisibilityHome, MkNoteVisibilitySpecif:
		s.Visibility = StatusVisibilityPublic
	case MkNoteVisibilityFollow:
		s.Visibility = StatusVisibilityPrivate
	default:
		s.Visibility = StatusVisibilityPrivate
	}
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
		re := n.ReNote.ToStatus(server)
		s.ReBlog = &re
	}
	return s
}

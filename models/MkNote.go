package models

import (
	"github.com/gizmo-ds/misstodon/internal/mfm"
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
	Reply        *MkNote          `json:"reply"`
	RepostID     *string          `json:"repostId"`
	Repost       *MkNote          `json:"repost"`
	ThreadId     *string          `json:"threadId"`
	Text         string           `json:"text"`
	Name         *string          `json:"name"`
	Cw           *string          `json:"cw"`
	UserId       string           `json:"userId"`
	User         *MkUser          `json:"user"`
	LocalOnly    bool             `json:"localOnly"`
	RenoteId     *string          `json:"renoteId"`
	RenoteCount  int              `json:"renoteCount"`
	RepliesCount int              `json:"repliesCount"`
	Reactions    map[string]int   `json:"reactions"`
	Visibility   MkNoteVisibility `json:"visibility"`
	Uri          *string          `json:"uri"`
	Url          *string          `json:"url"`
	Score        int              `json:"score"`
	FileIds      []string         `json:"fileIds"`
	Files        []MkFile         `json:"files"`
}

func (n *MkNote) ToStatus(server string) Status {
	s := Status{
		ID:        n.ID,
		Url:       "https://" + server + "/notes/" + n.ID,
		Uri:       "https://" + server + "/notes/" + n.ID,
		Tags:      []StatusTag{},
		CreatedAt: n.CreatedAt,
		Content:   n.Text,
		Emojis:    []struct{}{},
	}
	if content, err := mfm.ToHtml(n.Text, mfm.Option{Url: "https://" + server}); err == nil {
		s.Content = content
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
		s.MediaAttachments = append(s.MediaAttachments, file.ToMediaAttachment())
	}
	return s
}

package models

import (
	"time"

	"github.com/gizmo-ds/misstodon/internal/mfm"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/pkg/errors"
)

type MkUser struct {
	ID             string         `json:"id"`
	Username       string         `json:"username"`
	Name           string         `json:"name"`
	Host           *string        `json:"host,omitempty"`
	Location       *string        `json:"location"`
	Description    *string        `json:"description"`
	IsBot          bool           `json:"isBot"`
	IsLocked       bool           `json:"isLocked"`
	CreatedAt      string         `json:"createdAt,omitempty"`
	UpdatedAt      *string        `json:"updatedAt"`
	FollowersCount int            `json:"followersCount"`
	FollowingCount int            `json:"followingCount"`
	NotesCount     int            `json:"notesCount"`
	AvatarUrl      string         `json:"avatarUrl"`
	BannerUrl      string         `json:"bannerUrl"`
	Fields         []AccountField `json:"fields"`
	Instance       MkInstance     `json:"instance"`
	Mentions       []string       `json:"mentions"`
	IsMuted        bool           `json:"isMuted"`
	IsBlocked      bool           `json:"isBlocked"`
	IsBlocking     bool           `json:"isBlocking"`
	IsFollowing    bool           `json:"isFollowing"`
	IsFollowed     bool           `json:"isFollowed"`
}

type MkInstance struct {
	Name            string `json:"name"`
	SoftwareName    string `json:"softwareName"`
	SoftwareVersion string `json:"softwareVersion"`
	ThemeColor      string `json:"themeColor"`
	IconUrl         string `json:"iconUrl"`
	FaviconUrl      string `json:"faviconUrl"`
}

func (u *MkUser) ToAccount(ctx misstodon.Context) (Account, error) {
	var info Account
	var err error
	host := ctx.ProxyServer()
	if u.Host != nil {
		host = *u.Host
	}
	info = Account{
		ID:             u.ID,
		Username:       u.Username,
		Acct:           u.Username + "@" + host,
		DisplayName:    u.Name,
		Locked:         u.IsLocked,
		Bot:            u.IsBot,
		Url:            "https://" + host + "/@" + u.Username,
		Avatar:         u.AvatarUrl,
		AvatarStatic:   u.AvatarUrl,
		Header:         u.BannerUrl,
		HeaderStatic:   u.BannerUrl,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		StatusesCount:  u.NotesCount,
		Emojis:         []CustomEmoji{},
		Fields:         append([]AccountField{}, u.Fields...),
		CreatedAt:      u.CreatedAt,
		Limited:        &u.IsMuted,
	}
	if info.DisplayName == "" {
		info.DisplayName = info.Username
	}
	_lastStatusAt := u.UpdatedAt
	if _lastStatusAt != nil {
		lastStatusAt, err := time.Parse(time.RFC3339, *_lastStatusAt)
		if err != nil {
			return info, errors.WithStack(err)
		}
		t := lastStatusAt.Format("2006-01-02")
		info.LastStatusAt = &t
	}
	if u.Description != nil {
		info.Note, err = mfm.ToHtml(*u.Description, mfm.Option{
			Url:            utils.JoinURL(ctx.ProxyServer()),
			HashtagHandler: mfm.MastodonHashtagHandler,
		})
		if err != nil {
			return info, errors.WithStack(err)
		}
	}
	return info, nil
}

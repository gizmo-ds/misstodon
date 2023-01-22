package models

import (
	"time"

	"github.com/gizmo-ds/misstodon/internal/mfm"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/pkg/errors"
)

var (
	ErrAcctIsInvalid = errors.New("acct format is invalid")
)

type MkUser struct {
	ID             string         `json:"id"`
	Username       string         `json:"username"`
	Name           string         `json:"name"`
	Location       *string        `json:"location"`
	Description    *string        `json:"description"`
	IsBot          bool           `json:"isBot"`
	IsLocked       bool           `json:"isLocked"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      *string        `json:"updatedAt"`
	FollowersCount int            `json:"followersCount"`
	FollowingCount int            `json:"followingCount"`
	NotesCount     int            `json:"notesCount"`
	AvatarUrl      string         `json:"avatarUrl"`
	BannerUrl      string         `json:"bannerUrl"`
	Fields         []AccountField `json:"fields"`
	Instance       MkInstance     `json:"instance"`
}

type MkInstance struct {
	Name            string `json:"name"`
	SoftwareName    string `json:"softwareName"`
	SoftwareVersion string `json:"softwareVersion"`
	ThemeColor      string `json:"themeColor"`
	IconUrl         string `json:"iconUrl"`
	FaviconUrl      string `json:"faviconUrl"`
}

func (u *MkUser) ToAccount(acct, server string) (Account, error) {
	var info Account
	username, _host := utils.AcctInfo(acct)
	if username == "" {
		return info, errors.WithStack(ErrAcctIsInvalid)
	}
	if _host == "" {
		_host = server
	}
	createdAt, err := time.Parse(time.RFC3339, u.CreatedAt)
	if err != nil {
		err = errors.WithStack(err)
		return info, err
	}
	info = Account{
		ID:             u.ID,
		Username:       u.Username,
		Acct:           username + "@" + _host,
		DisplayName:    u.Name,
		Locked:         u.IsLocked,
		Bot:            u.IsBot,
		CreatedAt:      createdAt.Format("2006-01-02"),
		Url:            "https://" + _host + "/@" + username,
		Avatar:         u.AvatarUrl,
		AvatarStatic:   u.AvatarUrl,
		Header:         u.BannerUrl,
		HeaderStatic:   u.BannerUrl,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		StatusesCount:  u.NotesCount,
		Emojis:         []CustomEmoji{},
		Fields:         u.Fields,
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
		info.Note, err = mfm.ToHtml(*u.Description, mfm.Option{Url: "https://" + _host})
		if err != nil {
			return info, errors.WithStack(err)
		}
	}
	return info, nil
}

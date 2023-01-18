package models

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

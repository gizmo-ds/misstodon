package models

type Account struct {
	ID             string         `json:"id"`
	Username       string         `json:"username"`
	Acct           string         `json:"acct"`
	DisplayName    string         `json:"display_name"`
	Locked         bool           `json:"locked"`
	Bot            bool           `json:"bot"`
	Discoverable   bool           `json:"discoverable"`
	Group          bool           `json:"group"`
	CreatedAt      string         `json:"created_at"`
	Note           string         `json:"note"`
	Url            string         `json:"url"`
	Avatar         string         `json:"avatar"`
	AvatarStatic   string         `json:"avatar_static"`
	Header         string         `json:"header"`
	HeaderStatic   string         `json:"header_static"`
	FollowersCount int            `json:"followers_count"`
	FollowingCount int            `json:"following_count"`
	StatusesCount  int            `json:"statuses_count"`
	LastStatusAt   *string        `json:"last_status_at"`
	Emojis         []CustomEmoji  `json:"emojis"`
	Moved          *Account       `json:"moved,omitempty"`
	Suspended      *bool          `json:"suspended,omitempty"`
	Limited        *bool          `json:"limited,omitempty"`
	Fields         []AccountField `json:"fields"`
}

type AccountField struct {
	Name       string  `json:"name"`
	Value      string  `json:"value"`
	VerifiedAt *string `json:"verified_at,omitempty"`
}

func (a Account) SetLastStatusAt(t string) {
	a.LastStatusAt = &t
}

func (f *AccountField) SetVerifiedAt(t string) {
	f.VerifiedAt = &t
}

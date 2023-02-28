package models

type PostPrivacy = string

const (
	// PostPrivacyPublic Public post
	PostPrivacyPublic PostPrivacy = "public"
	// PostPrivacyUnlisted Unlisted post
	PostPrivacyUnlisted PostPrivacy = "unlisted"
	// PostPrivacyPrivate Followers-only post
	PostPrivacyPrivate PostPrivacy = "private"
	// PostPrivacyDirect Direct post
	PostPrivacyDirect PostPrivacy = "direct"
)

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

type CredentialAccount struct {
	Account
	Source struct {
		Privacy             PostPrivacy    `json:"privacy"`
		Sensitive           bool           `json:"sensitive"`
		Language            string         `json:"language"`
		Note                string         `json:"note"`
		Fields              []AccountField `json:"fields"`
		FollowRequestsCount int            `json:"follow_requests_count"`
	} `json:"source"`
	Role *struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Color       string `json:"color"`
		Position    int    `json:"position"`
		Permissions int    `json:"permissions"`
		Highlighted bool   `json:"highlighted"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	} `json:"role,omitempty"`
}

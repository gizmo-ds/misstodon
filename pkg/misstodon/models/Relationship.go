package models

type Relationship struct {
	ID         string `json:"id"`
	Following  bool   `json:"following"`
	FollowedBy bool   `json:"followed_by"`
	// ShowingReblogs bool     `json:"showing_reblogs"`
	Requested bool     `json:"requested"`
	Languages []string `json:"languages"`
	Blocking  bool     `json:"blocking"`
	BlockedBy bool     `json:"blocked_by"`
	Muting    bool     `json:"muting"`
	// MutingNotifications bool     `json:"muting_notifications"`
	// DomainBlocking bool   `json:"domain_blocking"`
	// Endorsed bool   `json:"endorsed"`
	// Notifying bool     `json:"notifying"`
	Note string `json:"note"`
}

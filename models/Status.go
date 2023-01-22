package models

type StatusVisibility string

const (
	// StatusVisibilityPublic Visible to everyone, shown in public timelines.
	StatusVisibilityPublic StatusVisibility = "public"
	// StatusVisibilityUnlisted Visible to public, but not included in public timelines.
	StatusVisibilityUnlisted StatusVisibility = "unlisted"
	// StatusVisibilityPrivate Visible to followers only, and to any mentioned users.
	StatusVisibilityPrivate StatusVisibility = "private"
	// StatusVisibilityDirect Visible only to mentioned users.
	StatusVisibilityDirect StatusVisibility = "direct"
)

type (
	Status struct {
		ID               string            `json:"id"`
		Uri              string            `json:"uri"`
		Url              string            `json:"url"`
		Visibility       StatusVisibility  `json:"visibility"`
		Tags             []StatusTag       `json:"tags"`
		CreatedAt        string            `json:"created_at"`
		EditedAt         *string           `json:"edited_at"`
		Content          string            `json:"content"`
		MediaAttachments []MediaAttachment `json:"media_attachments"`
		Card             *struct{}         `json:"card"`
		Emojis           []struct{}        `json:"emojis"`
	}
	StatusTag struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
)

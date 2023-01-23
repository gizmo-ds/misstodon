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
		ID                 string            `json:"id"`
		Uri                string            `json:"uri"`
		Url                string            `json:"url"`
		Visibility         StatusVisibility  `json:"visibility"`
		Tags               []StatusTag       `json:"tags"`
		CreatedAt          string            `json:"created_at"`
		EditedAt           *string           `json:"edited_at"`
		Content            string            `json:"content"`
		MediaAttachments   []MediaAttachment `json:"media_attachments"`
		Card               *struct{}         `json:"card"`
		Emojis             []struct{}        `json:"emojis"`
		Account            Account           `json:"account"`
		Sensitive          bool              `json:"sensitive"`
		SpoilerText        string            `json:"spoiler_text"`
		Bookmarked         bool              `json:"bookmarked"`
		Favourited         bool              `json:"favourited"`
		FavouritesCount    int               `json:"favourites_count"`
		InReplyToAccountId *string           `json:"in_reply_to_account_id"`
		InReplyToId        *string           `json:"in_reply_to_id"`
		Language           *string           `json:"language"`
		Mentions           []StatusMention   `json:"mentions"`
		Muted              bool              `json:"muted"`
		Poll               *struct{}         `json:"poll"`
		ReBlog             *Status           `json:"reblog"`
		ReBlogged          bool              `json:"reblogged"`
		ReBlogsCount       int               `json:"reblogs_count"`
		RepliesCount       int               `json:"replies_count"`
	}
	StatusTag struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	StatusMention struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Url      string `json:"url"`
		Acct     string `json:"acct"`
	}
)

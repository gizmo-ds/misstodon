package models

type (
	V1Instance struct {
		Uri              string         `json:"uri"`
		Title            string         `json:"title"`
		ShortDescription string         `json:"short_description"`
		Description      string         `json:"description"`
		Email            string         `json:"email"`
		Version          string         `json:"version"`
		Urls             V1InstanceUrls `json:"urls"`
		Stats            struct {
			UserCount   int `json:"user_count"`
			StatusCount int `json:"status_count"`
			DomainCount int `json:"domain_count"`
		} `json:"stats"`
		Thumbnail        string `json:"thumbnail"`
		Languages        any    `json:"languages"`
		Registrations    bool   `json:"registrations"`
		ApprovalRequired bool   `json:"approval_required"`
		InvitesEnabled   bool   `json:"invites_enabled"`
		Configuration    struct {
			Statuses struct {
				MaxCharacters            int `json:"max_characters"`
				MaxMediaAttachments      int `json:"max_media_attachments"`
				CharactersReservedPerUrl int `json:"characters_reserved_per_url"`
			} `json:"statuses"`
			MediaAttachments struct {
				SupportedMimeTypes []string `json:"supported_mime_types"`
			} `json:"media_attachments"`
		} `json:"configuration"`
		ContactAccount V1Account        `json:"contact_account"`
		Rules          []V1InstanceRule `json:"rules"`
	}
	V1InstanceUrls struct {
		StreamingApi string `json:"streaming_api"`
	}
	V1InstanceRule struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
)

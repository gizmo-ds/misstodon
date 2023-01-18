package models

type (
	Instance struct {
		Uri              string       `json:"uri"`
		Title            string       `json:"title"`
		ShortDescription string       `json:"short_description"`
		Description      string       `json:"description"`
		Email            string       `json:"email"`
		Version          string       `json:"version"`
		Urls             InstanceUrls `json:"urls"`
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
		ContactAccount Account        `json:"contact_account"`
		Rules          []InstanceRule `json:"rules"`
	}
	InstanceUrls struct {
		StreamingApi string `json:"streaming_api"`
	}
	InstanceRule struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
)

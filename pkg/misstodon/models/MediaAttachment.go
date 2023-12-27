package models

type MediaAttachment struct {
	ID          string  `json:"id"`
	BlurHash    string  `json:"blurhash"`
	Url         string  `json:"url"`
	Type        string  `json:"type"`
	TextUrl     *string `json:"text_url"`
	RemoteUrl   string  `json:"remote_url"`
	PreviewUrl  string  `json:"preview_url"`
	Description *string `json:"description"`
	Meta        struct {
		Small struct {
			Aspect float64 `json:"aspect"`
			Width  int     `json:"width"`
			Height int     `json:"height"`
			Size   string  `json:"size"`
		} `json:"small"`
		Original struct {
			Aspect float64 `json:"aspect"`
			Width  int     `json:"width"`
			Height int     `json:"height"`
			Size   string  `json:"size"`
		} `json:"original"`
	} `json:"meta"`
}

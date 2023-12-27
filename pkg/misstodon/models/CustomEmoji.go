package models

type CustomEmoji struct {
	Shortcode       string `json:"shortcode"`
	Url             string `json:"url"`
	StaticUrl       string `json:"static_url"`
	VisibleInPicker bool   `json:"visible_in_picker"`
	Category        string `json:"category"`
}

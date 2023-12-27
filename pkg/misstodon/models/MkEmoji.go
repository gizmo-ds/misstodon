package models

type MkEmoji struct {
	Aliases  []string `json:"aliases"`
	Name     string   `json:"name"`
	Category *string  `json:"category"`
	Url      string   `json:"url"`
}

func (e MkEmoji) ToCustomEmoji() CustomEmoji {
	r := CustomEmoji{
		Shortcode:       e.Name,
		Url:             e.Url,
		StaticUrl:       e.Url,
		VisibleInPicker: true,
	}
	if e.Category != nil {
		r.Category = *e.Category
	}
	return r
}

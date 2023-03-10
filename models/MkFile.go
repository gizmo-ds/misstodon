package models

import (
	"fmt"
	"strings"
)

type MkFile struct {
	ID           string  `json:"id"`
	ThumbnailUrl string  `json:"thumbnailUrl"`
	Type         string  `json:"type"`
	Url          string  `json:"url"`
	Name         string  `json:"name"`
	IsSensitive  bool    `json:"isSensitive"`
	Size         int64   `json:"size"`
	Md5          string  `json:"md5"`
	CreatedAt    string  `json:"createdAt"`
	BlurHash     *string `json:"blurhash"`
	Properties   struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"properties"`
}
type MkFolder struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"createdAt"`
	ParentId  *string `json:"parentId"`
}

func (f *MkFile) ToMediaAttachment() MediaAttachment {
	a := MediaAttachment{
		ID:         f.ID,
		Url:        f.Url,
		RemoteUrl:  f.Url,
		PreviewUrl: f.ThumbnailUrl,
	}
	a.Meta.Original.Width = f.Properties.Width
	a.Meta.Original.Height = f.Properties.Height
	if f.Properties.Width > 0 && f.Properties.Height > 0 {
		a.Meta.Original.Aspect = float64(f.Properties.Width) / float64(f.Properties.Height)
		a.Meta.Original.Size = fmt.Sprintf("%vx%v", f.Properties.Width, f.Properties.Height)
	}
	if f.BlurHash != nil {
		a.BlurHash = *f.BlurHash
	}
	t := strings.Split(f.Type, "/")[0]
	switch t {
	case "image":
		if f.Type == "image/gif" {
			a.Type = "gifv"
		} else {
			a.Type = "image"
		}
	case "application":
		if f.Type == "application/ogg" {
			a.Type = "audio"
		} else {
			a.Type = "unknown"
		}
	default:
		a.Type = t
	}
	return a
}

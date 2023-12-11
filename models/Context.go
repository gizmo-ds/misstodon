package models

type Context struct {
	Ancestors   []Status `json:"ancestors"`
	Descendants []Status `json:"descendants"`
}

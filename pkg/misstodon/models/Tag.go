package models

type Tag struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	History []struct {
		Day      string `json:"day"`
		Uses     string `json:"uses"`
		Accounts string `json:"accounts"`
	} `json:"history"`
	Following bool `json:"following,omitempty"`
}

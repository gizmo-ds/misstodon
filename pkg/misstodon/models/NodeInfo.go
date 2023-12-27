package models

type (
	NodeInfo struct {
		Version           string           `json:"version"`
		Software          NodeInfoSoftware `json:"software"`
		Protocols         []string         `json:"protocols"`
		Services          NodeInfoServices `json:"services"`
		Usage             NodeInfoUsage    `json:"usage"`
		OpenRegistrations bool             `json:"openRegistrations"`
		Metadata          any              `json:"metadata"`
	}
	NodeInfoSoftware struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	NodeInfoServices struct {
		Inbound  []string `json:"inbound"`
		Outbound []string `json:"outbound"`
	}
	NodeInfoUsage struct {
		Users struct {
			Total          int `json:"total"`
			ActiveMonth    int `json:"activeMonth"`
			ActiveHalfyear int `json:"activeHalfyear"`
		} `json:"users"`
		LocalPosts int `json:"localPosts"`
	}
)

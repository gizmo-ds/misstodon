package misskey

import (
	"github.com/gizmo-ds/misstodon/models"
)

func NodeInfo(server string, ni models.NodeInfo) (models.NodeInfo, error) {
	var result models.NodeInfo
	_, err := client.R().SetResult(&result).Get("https://" + server + "/nodeinfo/2.0")
	if err != nil {
		return ni, err
	}
	ni.Usage = result.Usage
	ni.OpenRegistrations = result.OpenRegistrations
	ni.Metadata = result.Metadata
	return ni, err
}

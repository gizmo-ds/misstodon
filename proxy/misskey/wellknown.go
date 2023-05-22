package misskey

import (
	"io"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
)

func NodeInfo(server string, ni models.NodeInfo) (models.NodeInfo, error) {
	var result models.NodeInfo
	_, err := client.R().
		SetResult(&result).
		Get(utils.JoinURL(server, "/nodeinfo/2.0"))
	if err != nil {
		return ni, err
	}
	ni.Usage = result.Usage
	ni.OpenRegistrations = result.OpenRegistrations
	ni.Metadata = result.Metadata
	return ni, err
}

func WebFinger(server, resource string, writer http.ResponseWriter) error {
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetQueryParam("resource", resource).
		Get(utils.JoinURL(server, "/.well-known/webfinger"))
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()
	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())
	_, err = io.Copy(writer, resp.RawBody())
	return err
}

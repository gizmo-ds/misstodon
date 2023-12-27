package misskey

import (
	"io"
	"net/http"

	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
)

func NodeInfo(server string, ni models.NodeInfo) (models.NodeInfo, error) {
	var result models.NodeInfo
	_, err := client.R().
		SetBaseURL(server).
		SetResult(&result).
		Get("/nodeinfo/2.0")
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
		SetBaseURL(server).
		SetDoNotParseResponse(true).
		SetQueryParam("resource", resource).
		Get("/.well-known/webfinger")
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()
	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())
	_, err = io.Copy(writer, resp.RawBody())
	return err
}

func HostMeta(server string, writer http.ResponseWriter) error {
	resp, err := client.R().
		SetBaseURL(server).
		SetDoNotParseResponse(true).
		Get("/.well-known/host-meta")
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()
	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())
	_, err = io.Copy(writer, resp.RawBody())
	return err
}

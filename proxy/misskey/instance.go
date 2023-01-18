package misskey

import (
	"net/http"
	"net/url"

	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

// SupportedMimeTypes is a list of supported mime types
//
// https://github.com/misskey-dev/misskey/blob/79212bbd375705f0fd658dd5b50b47f77d622fb8/packages/backend/src/const.ts#L25
var SupportedMimeTypes = []string{
	"image/png",
	"image/gif",
	"image/jpeg",
	"image/webp",
	"image/avif",
	"image/apng",
	"image/bmp",
	"image/tiff",
	"image/x-icon",
	"audio/opus",
	"video/ogg",
	"audio/ogg",
	"application/ogg",
	"video/quicktime",
	"video/mp4",
	"audio/mp4",
	"video/x-m4v",
	"audio/x-m4a",
	"video/3gpp",
	"video/3gpp2",
	"video/mpeg",
	"audio/mpeg",
	"video/webm",
	"audio/webm",
	"audio/aac",
	"audio/x-flac",
	"audio/vnd.wave",
}

func Instance(server, version string) (models.Instance, error) {
	var info models.Instance
	var serverInfo models.MkMeta
	resp, err := client.R().
		SetBody(map[string]any{
			"detail": false,
		}).
		SetResult(&serverInfo).
		Post("https://" + server + "/api/meta")
	serverUrl, err := url.Parse(serverInfo.URI)
	if err != nil {
		return info, err
	}
	if resp.StatusCode() != http.StatusOK {
		return info, errors.New("Failed to get instance info")
	}
	info = models.Instance{
		Uri:              serverUrl.Host,
		Title:            serverInfo.Name,
		Description:      serverInfo.Description,
		ShortDescription: serverInfo.Description,
		Email:            serverInfo.MaintainerEmail,
		Version:          version,
		Thumbnail:        serverInfo.BannerURL,
		Registrations:    !serverInfo.DisableRegistration,
		InvitesEnabled:   serverInfo.Policies.CanInvite,
		Rules:            []models.InstanceRule{},
		Languages:        serverInfo.Langs,
	}
	//TODO: 需要先实现 `/streaming`
	//info.Urls.StreamingApi = serverInfo.StreamingAPI
	if info.Languages == nil {
		info.Languages = []string{}
	}
	info.Configuration.Statuses.MaxCharacters = serverInfo.MaxNoteTextLength
	// NOTE: misskey没有相关限制, 此处返回固定值
	info.Configuration.Statuses.MaxMediaAttachments = 4
	// NOTE: misskey没有相关设置, 此处返回固定值
	info.Configuration.Statuses.CharactersReservedPerUrl = 23
	info.Configuration.MediaAttachments.SupportedMimeTypes = SupportedMimeTypes

	var serverStats models.MkStats
	resp, err = client.R().
		SetBody(map[string]any{}).
		SetResult(&serverStats).
		Post("https://" + server + "/api/stats")
	if err != nil {
		return info, err
	}
	if resp.StatusCode() != http.StatusOK {
		return info, errors.New("Failed to get instance info")
	}
	info.Stats.UserCount = serverStats.OriginalUsersCount
	info.Stats.StatusCount = serverStats.OriginalNotesCount
	info.Stats.DomainCount = serverStats.Instances
	return info, err
}

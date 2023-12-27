package misskey

import (
	"io"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/pkg/errors"
)

func driveFileCreate(ctx misstodon.Context, filename string, content io.Reader) (models.MkFile, error) {
	var file models.MkFile
	// find folder
	folders, err := driveFolders(ctx)
	if err != nil {
		return file, err
	}
	var saveFolder *models.MkFolder
	for _, folder := range folders {
		if folder.Name == "misstodon" {
			saveFolder = &folder
			break
		}
	}

	// create folder if not exists
	if saveFolder == nil {
		folder, err := driveFolderCreate(ctx, "misstodon")
		if err != nil {
			return file, err
		}
		saveFolder = &folder
	}

	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetFormData(map[string]string{
			"folderId":    saveFolder.Id,
			"name":        filename,
			"i":           *ctx.Token(),
			"force":       "true",
			"isSensitive": "false",
		}).
		SetMultipartField("file", filename, "application/octet-stream", content).
		SetResult(&file).
		Post("/api/drive/files/create")
	if err != nil {
		return file, err
	}
	if resp.StatusCode() != http.StatusOK {
		return file, errors.New("failed to verify credentials")
	}
	return file, nil
}

func driveFolders(ctx misstodon.Context) (folders []models.MkFolder, err error) {
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(utils.Map{"i": ctx.Token(), "limit": 100}).
		SetResult(&folders).
		Post("/api/drive/folders")
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		return folders, errors.New("failed to verify credentials")
	}
	return
}

func driveFolderCreate(ctx misstodon.Context, name string) (models.MkFolder, error) {
	var folder models.MkFolder
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(utils.Map{"name": name, "i": ctx.Token()}).
		SetResult(&folder).
		Post("/api/drive/folders/create")
	if err != nil {
		return folder, err
	}
	if resp.StatusCode() != http.StatusOK {
		return folder, errors.New("failed to verify credentials")
	}
	return folder, nil
}

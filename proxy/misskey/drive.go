package misskey

import (
	"io"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func driveFileCreate(ctx Context, filename string, content io.Reader) (models.MkFile, error) {
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
		SetFormData(map[string]string{
			"folderId":    saveFolder.Id,
			"name":        filename,
			"i":           *ctx.Token(),
			"force":       "true",
			"isSensitive": "false",
		}).
		SetMultipartField("file", filename, "application/octet-stream", content).
		SetResult(&file).
		Post(utils.JoinURL(ctx.Server(), "/api/drive/files/create"))
	if err != nil {
		return file, err
	}
	if resp.StatusCode() != 200 {
		return file, errors.New("failed to verify credentials")
	}
	return file, nil
}

func driveFolders(ctx Context) (folders []models.MkFolder, err error) {
	resp, err := client.R().
		SetBody(utils.Map{"i": ctx.Token(), "limit": 100}).
		SetResult(&folders).
		Post(utils.JoinURL(ctx.Server(), "/api/drive/folders"))
	if err != nil {
		return
	}
	if resp.StatusCode() != 200 {
		return folders, errors.New("failed to verify credentials")
	}
	return
}

func driveFolderCreate(ctx Context, name string) (models.MkFolder, error) {
	var folder models.MkFolder
	resp, err := client.R().
		SetBody(utils.Map{"name": name, "i": ctx.Token()}).
		SetResult(&folder).
		Post(utils.JoinURL(ctx.Server(), "/api/drive/folders/create"))
	if err != nil {
		return folder, err
	}
	if resp.StatusCode() != 200 {
		return folder, errors.New("failed to verify credentials")
	}
	return folder, nil
}

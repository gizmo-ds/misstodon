package misskey

import (
	"mime/multipart"

	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func MediaUpload(server, token string, file *multipart.FileHeader, description string) (models.MediaAttachment, error) {
	var ma models.MediaAttachment
	if file == nil {
		return ma, errors.New("file is nil")
	}
	f, err := file.Open()
	if err != nil {
		return ma, err
	}
	defer f.Close()

	fileInfo, err := driveFileCreate(server, token, file.Filename, f)
	if err != nil {
		return ma, err
	}
	ma = fileInfo.ToMediaAttachment()
	return ma, nil
}

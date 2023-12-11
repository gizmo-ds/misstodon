package misskey

import (
	"mime/multipart"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func MediaUpload(ctx misstodon.Context, file *multipart.FileHeader, description string) (models.MediaAttachment, error) {
	var ma models.MediaAttachment
	if file == nil {
		return ma, errors.New("file is nil")
	}
	f, err := file.Open()
	if err != nil {
		return ma, err
	}
	defer f.Close()

	fileInfo, err := driveFileCreate(ctx, file.Filename, f)
	if err != nil {
		return ma, err
	}
	ma = fileInfo.ToMediaAttachment()
	return ma, nil
}

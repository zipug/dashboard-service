package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/pkg/mime"
	"fmt"
)

func (d *DashboardService) UploadAttachment(name, ext string, content []byte, user_id int64) (string, error) {
	ctx := context.Background()
	d.log.Log("info", "uploading attachment", logger.WithStrAttr("name", name))
	contentType, err := mime.ConvertExtToMIME(ext)
	if err != nil {
		d.log.Log("error", "error while converting extension to MIME", logger.WithErrAttr(err))
		return "", err
	}
	resp, err := d.minio.UploadFile(
		ctx,
		models.File{
			Name:        name,
			Data:        content,
			Bucket:      "attachments",
			ContentType: contentType,
		},
	)
	if err != nil {
		d.log.Log("error", "error while uploading attachment", logger.WithErrAttr(err))
		return "", err
	}
	d.log.Log("info", "attachment successfully uploaded", logger.WithStrAttr("url", resp.Url))
	d.attachment.CreateAttachment(ctx, models.Attachment{
		Name:   fmt.Sprintf("%s.%s", name, ext),
		UserId: user_id,
	})
	return resp.Url, nil
}

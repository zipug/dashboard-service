package application

import (
	"context"
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/pkg/mime"
	"errors"
)

func (d *DashboardService) UploadAttachment(name, ext string, content []byte, user_id int64) (dto.AttachmentDto, error) {
	ctx := context.Background()
	d.log.Log("info", "uploading attachment", logger.WithStrAttr("name", name))
	contentType, err := mime.ConvertExtToMIME(ext)
	if err != nil {
		d.log.Log("error", "error while converting extension to MIME", logger.WithErrAttr(err))
		return dto.AttachmentDto{}, err
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
		return dto.AttachmentDto{}, err
	}
	d.log.Log("info", "attachment successfully uploaded", logger.WithStrAttr("url", resp.Url))
	db_attachment, err := d.attachment.CreateAttachment(ctx, models.Attachment{
		Name:     name,
		UserId:   user_id,
		Mimetype: contentType,
		ObjectId: resp.ObjectId,
	})
	if err != nil {
		d.log.Log("error", "error while creating attachment", logger.WithErrAttr(err))
		return dto.AttachmentDto{}, err
	}
	if db_attachment.Id == -1 {
		d.log.Log("error", "error while creating attachment", logger.WithErrAttr(err))
		return dto.AttachmentDto{}, errors.New("error while creating attachment")
	}
	resp_dto := dto.ToAttachmentDto(db_attachment)
	resp_dto.URL = resp.Url
	return resp_dto, nil
}

func (d *DashboardService) BindAttachment(attachment_id, article_id, user_id int64) error {
	ctx := context.Background()
	err := d.attachment.BindAttachment(ctx, attachment_id, article_id, user_id)
	if err != nil {
		d.log.Log("error", "error while binding attachment", logger.WithErrAttr(err))
		return err
	}
	d.log.Log(
		"info",
		"attachment successfully binded",
		logger.WithInt64Attr("attachment_id", attachment_id),
		logger.WithInt64Attr("article_id", article_id),
	)
	return nil
}

func (d *DashboardService) GetAttachmentById(attacment_id, user_id int64) (models.Attachment, error) {
	ctx := context.Background()
	attachment, err := d.attachment.GetAttachmentById(ctx, attacment_id, user_id)
	if err != nil {
		d.log.Log("error", "error while getting attachment", logger.WithErrAttr(err))
		return models.Attachment{}, err
	}
	data, err := d.minio.GetFileUrl(ctx, attachment.ObjectId, "attachments")
	if err != nil {
		d.log.Log("error", "error while getting attachment", logger.WithErrAttr(err))
		return models.Attachment{}, err
	}
	attachment.URL = data.Url
	d.log.Log("info", "successfully got attachment", logger.WithInt64Attr("attachment_id", attacment_id))
	return attachment, nil
}

func (d *DashboardService) GetAllAttachments(user_id int64) ([]models.Attachment, error) {
	ctx := context.Background()
	attachments, err := d.attachment.GetAllAttachments(ctx, user_id)
	if err != nil {
		d.log.Log("error", "error while getting attachment", logger.WithErrAttr(err))
		return nil, err
	}
	objectIds := make([]string, len(attachments))
	for i, val := range attachments {
		objectIds[i] = val.ObjectId
	}
	data, err := d.minio.GetManyFileUrls(ctx, objectIds, "attachments")
	if err != nil {
		d.log.Log("error", "error while getting attachment", logger.WithErrAttr(err))
		return nil, err
	}
	for i := 0; i < len(attachments); i++ {
		if a, ok := data[attachments[i].ObjectId]; ok {
			attachments[i].URL = a.Url
		}
	}
	d.log.Log("info", "successfully got all attachments")
	return attachments, nil
}

func (d *DashboardService) DeleteAttachment(attachment_id, user_id int64) error {
	ctx := context.Background()
	att, err := d.attachment.GetAttachmentById(ctx, attachment_id, user_id)
	if err != nil {
		return err
	}
	if err := d.attachment.DeleteAttachment(ctx, attachment_id, user_id); err != nil {
		d.log.Log("error", "error while deleting attachment", logger.WithErrAttr(err))
		return err
	}
	if err = d.minio.DeleteFile(ctx, att.ObjectId, "attachments"); err != nil {
		d.log.Log("error", "error while deleting attachment", logger.WithErrAttr(err))
		return err
	}
	d.log.Log(
		"info",
		"attachment successfully deleted",
		logger.WithInt64Attr("attachment_id", attachment_id),
	)
	return nil
}

package attachments

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type AttachmentsService struct {
	repo ports.AttachmentsRepository
}

func NewAttachmentsService(repo ports.AttachmentsRepository) *AttachmentsService {
	return &AttachmentsService{repo: repo}
}

func (a *AttachmentsService) GetAttachmentById(ctx context.Context, attachment_id, user_id int64) (models.Attachment, error) {
	attachment, err := a.repo.GetAttachmentById(ctx, attachment_id, user_id)
	if err != nil {
		return models.Attachment{}, err
	}
	return attachment.ToValue(), nil
}

func (a *AttachmentsService) GetAllAttachments(ctx context.Context, user_id int64) ([]models.Attachment, error) {
	attachments, err := a.repo.GetAllAttachments(ctx, user_id)
	if err != nil {
		return nil, err
	}
	var res []models.Attachment
	for _, val := range attachments {
		res = append(res, val.ToValue())
	}
	return res, nil
}

func (a *AttachmentsService) CreateAttachment(ctx context.Context, attachment models.Attachment) (models.Attachment, error) {
	dbo := dto.ToAttachmentDbo(attachment)
	new_dbo, err := a.repo.CreateAttachment(ctx, dbo)
	if err != nil {
		return models.Attachment{}, err
	}
	return new_dbo.ToValue(), nil
}

func (a *AttachmentsService) BindAttachment(ctx context.Context, attachment_id, article_id, user_id int64) error {
	return a.repo.BindAttachment(ctx, attachment_id, article_id, user_id)
}

func (a *AttachmentsService) DeleteAttachment(ctx context.Context, attachment_id, user_id int64) error {
	return nil
}

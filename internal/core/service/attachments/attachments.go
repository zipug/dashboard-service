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

func (a *AttachmentsService) GetAttachmentById(ctx context.Context, object_id string, user_id int64) (models.Attachment, error) {
	return models.Attachment{}, nil
}

func (a *AttachmentsService) GetAllAttachments(ctx context.Context, user_id int64) ([]models.Attachment, error) {
	return []models.Attachment{}, nil
}

func (a *AttachmentsService) CreateAttachment(ctx context.Context, attachment models.Attachment) (int64, error) {
	dbo := dto.ToAttachmentDbo(attachment)
	return a.repo.CreateAttachment(ctx, dbo)
}

func (a *AttachmentsService) BindAttachment(ctx context.Context, object_id string, article_id, user_id int64) error {
	return nil
}

func (a *AttachmentsService) DeleteAttachment(ctx context.Context, object_id string, user_id int64) error {
	return nil
}

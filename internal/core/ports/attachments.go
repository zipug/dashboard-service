package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type AttachmentsService interface {
	GetAttachmentById(context.Context, int64, int64) (models.Attachment, error)
	GetAllAttachments(context.Context, int64) ([]models.Attachment, error)
	CreateAttachment(context.Context, models.Attachment) (models.Attachment, error)
	BindAttachment(context.Context, int64, int64, int64) error
	DeleteAttachment(context.Context, int64, int64) error
}

type AttachmentsRepository interface {
	GetAttachmentById(context.Context, int64, int64) (*dto.AttachmentDbo, error)
	GetAllAttachments(context.Context, int64) ([]dto.AttachmentDbo, error)
	CreateAttachment(context.Context, dto.AttachmentDbo) (dto.AttachmentDbo, error)
	BindAttachment(context.Context, int64, int64, int64) error
	DeleteAttachment(context.Context, int64, int64) error
}

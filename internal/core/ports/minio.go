package ports

import (
	"context"
	"dashboard/internal/core/models"
)

type MinioService interface {
	UploadFile(context.Context, models.File) (string, error)
	UploadManyFiles(context.Context, []models.File) ([]string, error)
	GetFileUrl(context.Context, string, string) (string, error)
	GetManyFileUrls(context.Context, []string, string) ([]string, error)
	DeleteFile(context.Context, string, string) error
	DeleteManyFiles(context.Context, []string, string) error
}

type MinioRepository interface {
	UploadFile(context.Context, models.File) (string, error)
	UploadManyFiles(context.Context, []models.File) ([]string, []models.MinioErr)
	GetFileUrl(context.Context, string, string) (string, error)
	GetManyFileUrls(context.Context, []string, string) ([]string, []models.MinioErr)
	DeleteFile(context.Context, string, string) error
	DeleteManyFiles(context.Context, []string, string) []models.MinioErr
}

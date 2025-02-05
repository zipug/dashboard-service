package minio

import (
	"context"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"errors"
	"fmt"
	"strings"
)

type MinioService struct {
	repo ports.MinioRepository
}

func NewMinioService(repo ports.MinioRepository) *MinioService {
	return &MinioService{repo: repo}
}

func (m *MinioService) UploadFile(ctx context.Context, file models.File) (models.MinioResponse, error) {
	return m.repo.UploadFile(ctx, file)
}

func (m *MinioService) UploadManyFiles(ctx context.Context, files []models.File) ([]models.MinioResponse, error) {
	urls, errs := m.repo.UploadManyFiles(ctx, files)
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return urls, errors.New(strings.Join(err_array, "\n"))
}

func (m *MinioService) GetFileUrl(ctx context.Context, object_id, bucket string) (models.MinioResponse, error) {
	return m.repo.GetFileUrl(ctx, object_id, bucket)
}

func (m *MinioService) GetManyFileUrls(ctx context.Context, object_ids []string, bucket string) ([]models.MinioResponse, error) {
	urls, errs := m.repo.GetManyFileUrls(ctx, object_ids, bucket)
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return urls, errors.New(strings.Join(err_array, "\n"))
}

func (m *MinioService) DeleteFile(ctx context.Context, object_id, bucket string) error {
	return m.repo.DeleteFile(ctx, object_id, bucket)
}

func (m *MinioService) DeleteManyFiles(ctx context.Context, object_ids []string, bucket string) error {
	errs := m.repo.DeleteManyFiles(ctx, object_ids, bucket)
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return errors.New(strings.Join(err_array, "\n"))
}

package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type GeneratedReportsService interface {
	GetAllGeneratedReports(context.Context, int64) ([]models.GeneratedReport, error)
	GetGeneratedReportById(context.Context, int64, int64) (models.GeneratedReport, error)
	DeleteGeneratedReport(context.Context, int64, int64) error
}

type GeneratedReportsRepository interface {
	GetAllGeneratedReports(context.Context, int64) ([]dto.GeneratedReportDbo, error)
	GetGeneratedReportById(context.Context, int64, int64) (*dto.GeneratedReportDbo, error)
	DeleteGeneratedReport(context.Context, int64, int64) error
}

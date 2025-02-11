package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type ReportsService interface {
	GetAllReports(context.Context) ([]models.Report, error)
	ExecuteReport(context.Context, int64, int64, string, string) error /* ctx, report_id, user_id, date_from, date_to */
}

type ReportsPostgresRepository interface {
	GetAllReports(context.Context) ([]dto.ReportDbo, error)
}

type ReportsRedisRepository interface {
	ExecuteReport(context.Context, int64, int64, string, string) error
}

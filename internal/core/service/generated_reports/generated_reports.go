package generatedreports

import (
	"context"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type GeneratedReportsService struct {
	repo ports.GeneratedReportsRepository
}

func NewGeneratedReportsService(repo ports.GeneratedReportsRepository) *GeneratedReportsService {
	return &GeneratedReportsService{repo: repo}
}

func (s *GeneratedReportsService) GetAllGeneratedReports(ctx context.Context, user_id int64) ([]models.GeneratedReport, error) {
	dbo, err := s.repo.GetAllGeneratedReports(ctx, user_id)
	if err != nil {
		return nil, err
	}
	var res []models.GeneratedReport
	for _, r := range dbo {
		res = append(res, r.ToNullContentValue())
	}
	return res, nil
}

func (s *GeneratedReportsService) GetGeneratedReportById(ctx context.Context, generated_report_id, user_id int64) (models.GeneratedReport, error) {
	dbo, err := s.repo.GetGeneratedReportById(ctx, generated_report_id, user_id)
	if err != nil {
		return models.GeneratedReport{}, err
	}
	return dbo.ToValue()
}

func (s *GeneratedReportsService) DeleteGeneratedReport(ctx context.Context, generated_report_id, user_id int64) error {
	return s.repo.DeleteGeneratedReport(ctx, generated_report_id, user_id)
}

package reports

import (
	"context"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type ReportsService struct {
	db    ports.ReportsPostgresRepository
	redis ports.ReportsRedisRepository
}

func NewReportsService(db ports.ReportsPostgresRepository, redis ports.ReportsRedisRepository) *ReportsService {
	return &ReportsService{db: db, redis: redis}
}

func (s *ReportsService) GetAllReports(ctx context.Context) ([]models.Report, error) {
	dbo, err := s.db.GetAllReports(ctx)
	if err != nil {
		return nil, err
	}
	var res []models.Report
	for _, r := range dbo {
		res = append(res, r.ToValue())
	}
	return res, nil
}

func (s *ReportsService) ExecuteReport(ctx context.Context, report_id, user_id int64, date_from, date_to string) error {
	return s.redis.ExecuteReport(ctx, report_id, user_id, date_from, date_to)
}

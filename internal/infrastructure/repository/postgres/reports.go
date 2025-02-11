package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var ErrReportsNotFound = errors.New("could not find any reports")

func (repo *PostgresRepository) GetAllReports(ctx context.Context) ([]dto.ReportDbo, error) {
	rows, err := pu.Dispatch[dto.ReportDbo](
		ctx,
		repo.db,
		`
    SELECT r.id, r.name, r.description, r.icon
		FROM reports r
		WHERE r.deleted_at IS NULL;
		`,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrReportsNotFound
	}
	return rows, nil
}

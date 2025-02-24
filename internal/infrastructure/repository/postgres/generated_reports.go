package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrGenRepNotFound   = errors.New("could not find generated_report by id")
	ErrGenRepsNotFound  = errors.New("could not find any generated_reports")
	ErrGenRepNotDeleted = errors.New("could not delete generated_report")
)

func (repo *PostgresRepository) GetGeneratedReportById(ctx context.Context, generated_report_id, user_id int64) (*dto.GeneratedReportDbo, error) {
	rows, err := pu.Dispatch[dto.GeneratedReportDbo](
		ctx,
		repo.db,
		`
		SELECT g.id, g.user_id, g.report_id, g.object_id, g.content::text, g.date_from, g.date_to
		FROM generated_reports g
		LEFT JOIN user_roles ur ON ur.user_id = $2::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE g.id = $1::bigint
		  AND (g.user_id = $2::bigint OR r.name = 'admin')
		  AND g.deleted_at IS NULL;
		`,
		generated_report_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGenRepNotFound
	}
	return &rows[0], nil
}

func (repo *PostgresRepository) GetAllGeneratedReports(ctx context.Context, user_id int64) ([]dto.GeneratedReportDbo, error) {
	rows, err := pu.Dispatch[dto.GeneratedReportDbo](
		ctx,
		repo.db,
		`
		SELECT g.id, g.user_id, g.report_id, g.object_id, g.content::text, g.date_from, g.date_to
		FROM generated_reports g
		LEFT JOIN user_roles ur ON ur.user_id = $1::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE (g.user_id = $1::bigint OR r.name = 'admin')
		  AND g.deleted_at IS NULL;
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGenRepsNotFound
	}
	return rows, nil
}

func (repo *PostgresRepository) DeleteGeneratedReport(ctx context.Context, generated_report_id, user_id int64) error {
	_, err := pu.Dispatch[dto.GeneratedReportDbo](
		ctx,
		repo.db,
		`
		DELETE FROM generated_reports
		USING generated_reports AS g
		LEFT JOIN user_roles ur ON ur.user_id = $2::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE generated_reports.id = $1::bigint
		  AND (g.user_id = $2::bigint OR r.name = 'admin');
		`,
		generated_report_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

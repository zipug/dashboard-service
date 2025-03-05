package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var ErrCreateDialog = errors.New("could not create dialog")

func (repo *PostgresRepository) CreateDialog(
	ctx context.Context,
	user_id,
	statistics_id int64,
	answer string,
) error {
	tx := repo.db.MustBegin()
	rows, err := pu.DispatchTx[dto.StatisticDbo](
		ctx,
		tx,
		`
		INSERT INTO statistics (bot_id, telegram_id, question, article_name, is_resolved, parent_id)
		SELECT s.bot_id, s.telegram_id, $1::text, 'support', true, s.id
		FROM statistics s
		WHERE s.id = $2::bigint
		RETURNING *;
		`,
		answer,
		statistics_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(rows) == 0 {
		tx.Rollback()
		return ErrCreateDialog
	}
	upd_rows, err := pu.DispatchTx[dto.StatisticDbo](
		ctx,
		tx,
		`
		UPDATE statistics
		SET is_resolved = TRUE
		WHERE id = $1::bigint AND is_resolved is FALSE
		RETURNING *;
		`,
		statistics_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(upd_rows) == 0 {
		tx.Rollback()
		return ErrCreateDialog
	}
	tx.Commit()
	return nil
}

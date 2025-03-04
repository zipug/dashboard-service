package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrChatsNotFound   = errors.New("chats not found")
	ErrChatNotFound    = errors.New("chat not found")
	ErrResolveQuestion = errors.New("error while resolving question")
)

func (repo *PostgresRepository) GetAllChats(ctx context.Context, user_id int64) ([]dto.ChatDbo, error) {
	rows, err := pu.Dispatch[dto.ChatDbo](
		ctx,
		repo.db,
		`
		SELECT s.id,
		       s.bot_id,
		       s.telegram_id,
		       b.project_id,
		       s.question,
		       s.created_at,
		       p.name,
		       u.created_by,
		       u.id AS user_id,
		       s.is_resolved
		FROM statistics s
			LEFT JOIN bots b ON b.id = s.bot_id
			LEFT JOIN projects p ON p.id = b.project_id
			LEFT JOIN users u ON p.user_id = u.created_by
		WHERE (u.id = $1::bigint OR u.created_by = $1::bigint);
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrChatsNotFound
	}
	return rows, nil
}

func (repo *PostgresRepository) GetChatById(ctx context.Context, chat_id, user_id int64) (*dto.ChatDbo, error) {
	rows, err := pu.Dispatch[dto.ChatDbo](
		ctx,
		repo.db,
		`
		SELECT s.id,
		       s.bot_id,
		       s.telegram_id,
		       b.project_id,
		       s.question,
		       s.created_at,
		       p.name,
		       u.created_by,
		       u.id AS user_id,
		       s.is_resolved
		FROM statistics s
			LEFT JOIN bots b ON b.id = s.bot_id
			LEFT JOIN projects p ON p.id = b.project_id
			LEFT JOIN users u ON p.user_id = u.created_by
		WHERE s.id = $1::bigint
		  AND (u.id = $2::bigint OR u.created_by = $2::bigint);
		`,
		chat_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrChatNotFound
	}
	row := rows[0]
	return &row, nil
}

func (repo *PostgresRepository) ResolveQuestion(ctx context.Context, chat_id int64) error {
	rows, err := pu.Dispatch[dto.StatisticDbo](
		ctx,
		repo.db,
		`
		UPDATE statistics
		SET is_resolved = TRUE
		WHERE id = $1::bigint AND is_resolved is FALSE
		RETURNING *;
		`,
		chat_id,
	)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return ErrResolveQuestion
	}
	return nil
}

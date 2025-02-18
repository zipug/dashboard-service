package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrBotNotFound   = errors.New("could not find bot by id")
	ErrBotsNBotFound = errors.New("could not find any bots")
	ErrBotNotCreated = errors.New("could not create bot")
	ErrBotNotUpdated = errors.New("could not update bot")
	ErrBotNotDeleted = errors.New("could not delete bot")
)

func (repo *PostgresRepository) GetBotById(ctx context.Context, bot_id, user_id int64) (*dto.BotDbo, error) {
	rows, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
		SELECT b.id, b.project_id, b.name, b.description, b.icon, b.state, b.user_id, b.api_token
		FROM bots b
		WHERE b.id = $1::bigint
		  AND b.user_id = $2::bigint
		  AND b.deleted_at IS NULL;
		`,
		bot_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrBotNotFound
	}
	return &rows[0], nil
}

func (repo *PostgresRepository) GetAllBots(ctx context.Context, user_id int64) ([]dto.BotDbo, error) {
	rows, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
		SELECT b.id, b.project_id, b.name, b.description, b.icon, b.state, b.user_id, b.api_token
		FROM bots b
		WHERE b.user_id = $1::bigint
		  AND b.deleted_at IS NULL;
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrBotsNBotFound
	}
	return rows, nil
}

func (repo *PostgresRepository) CreateBot(ctx context.Context, bot dto.BotDbo, user_id int64) (int64, error) {
	rows, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
		INSERT INTO bots (project_id, name, description, icon, state, user_id, api_token)
		VALUES ($1::bigint, $2::text, $3::text, $4::text, $5::text, $6::bigint, $7::text)
		RETURNING *;
		`,
		bot.ProjectId,
		bot.Name,
		bot.Description,
		bot.Icon,
		bot.State,
		user_id,
		bot.ApiToken,
	)
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return -1, ErrBotNotCreated
	}
	return rows[0].Id, nil
}

func (repo *PostgresRepository) UpdateBotById(ctx context.Context, bot dto.BotDbo, user_id int64) (*dto.BotDbo, error) {
	rows, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
    UPDATE bots
		SET name = $1::text,
		    description = $2::text,
		    icon = $3::text,
		    api_token = $4::text
		WHERE id = $5::bigint
		  AND user_id = $6::bigint
		RETURNING *;
		`,
		bot.Name,
		bot.Description,
		bot.Icon,
		bot.ApiToken,
		bot.Id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrBotNotUpdated
	}
	return &rows[0], nil
}

func (repo *PostgresRepository) DeleteBotById(ctx context.Context, bot_id, user_id int64) error {
	_, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
		DELETE FROM bots
		USING bots AS b
		LEFT JOIN user_roles ur ON b.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE bots.id = $1::bigint
		  AND bots.state != 'running'
		  AND (b.user_id = $2::bigint OR r.name = 'admin');
		`,
		bot_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) SetBotState(ctx context.Context, state string, bot_id, user_id int64) error {
	_, err := pu.Dispatch[dto.BotDbo](
		ctx,
		repo.db,
		`
		UPDATE bots
		SET state = $1::text
		WHERE id = $2::bigint
		  AND user_id = $3::bigint;
		`,
		state,
		bot_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

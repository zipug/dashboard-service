package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var ErrCreateAttachment = errors.New("could not create attachment")

func (repo *PostgresRepository) GetAttachmentById(ctx context.Context, object_id string, user_id int64) (*dto.AttachmentDbo, error) {
	return nil, nil
}

func (repo *PostgresRepository) GetAllAttachments(ctx context.Context, user_id int64) ([]dto.AttachmentDbo, error) {
	return nil, nil
}

func (repo *PostgresRepository) CreateAttachment(ctx context.Context, attachment dto.AttachmentDbo) (int64, error) {
	rows, err := pu.Dispatch[dto.AttachmentDbo](
		ctx,
		repo.db,
		`
			INSERT INTO attachments(name, description, mimetype, object_id, user_id)
			VAlUES($1::text, $2::text, $3::text, $4::text, $5::bigint)
			RETURNING *;
		`,
		attachment.Name,
		attachment.Description,
		attachment.Mimetype,
		attachment.ObjectId,
		attachment.UserID,
	)
	if err != nil {
		return -1, err
	}
	if len(rows) == 0 {
		return -1, ErrCreateAttachment
	}
	return rows[0].Id, err
}

func (repo *PostgresRepository) BindAttachment(ctx context.Context, object_id string, article_id, user_id int64) error {
	return nil
}

func (repo *PostgresRepository) DeleteAttachment(ctx context.Context, object_id string, user_id int64) error {
	return nil
}

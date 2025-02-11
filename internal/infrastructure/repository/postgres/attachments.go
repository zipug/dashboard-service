package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrCreateAttachment  = errors.New("could not create attachment")
	ErrBindAttachment    = errors.New("could not bind attachment")
	ErrGetAttachment     = errors.New("could not get attachment")
	ErrGetAllAttachments = errors.New("could not get all attachments")
	ErrDeleteAttachment  = errors.New("could not delete attachment")
)

func (repo *PostgresRepository) GetAttachmentById(ctx context.Context, attachment_id, user_id int64) (*dto.AttachmentDbo, error) {
	rows, err := pu.Dispatch[dto.AttachmentDbo](
		ctx,
		repo.db,
		`
		SELECT id, name, description, object_id, mimetype, user_id
		FROM attachments
		WHERE id = $1::bigint
		  AND user_id = $2::bigint
		  AND deleted_at IS NULL;
		`,
		attachment_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGetAttachment
	}
	return &rows[0], nil
}

func (repo *PostgresRepository) GetAllAttachments(ctx context.Context, user_id int64) ([]dto.AttachmentDbo, error) {
	rows, err := pu.Dispatch[dto.AttachmentDbo](
		ctx,
		repo.db,
		`
		SELECT id, name, description, object_id, mimetype, user_id
		FROM attachments
		WHERE user_id = $1::bigint
		  AND deleted_at IS NULL;
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGetAllAttachments
	}
	return rows, nil
}

func (repo *PostgresRepository) GetAllAttachmentsByArticleId(ctx context.Context, article_id int64) ([]dto.AttachmentDbo, error) {
	rows, err := pu.Dispatch[dto.AttachmentDbo](
		ctx,
		repo.db,
		`
		SELECT a.id, a.name, a.description, a.object_id, a.mimetype, a.user_id
    FROM attachments a
		LEFT JOIN attachments_articles aa ON a.id = aa.attachment_id
		WHERE aa.article_id = $1::bigint;
		`,
		article_id,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *PostgresRepository) CreateAttachment(ctx context.Context, attachment dto.AttachmentDbo) (dto.AttachmentDbo, error) {
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
		return dto.AttachmentDbo{}, err
	}
	if len(rows) == 0 {
		return dto.AttachmentDbo{}, ErrCreateAttachment
	}
	return rows[0], err
}

func (repo *PostgresRepository) BindAttachment(ctx context.Context, attachment_id, article_id, user_id int64) error {
	rows, err := pu.Dispatch[dto.AttachmentArticleDbo](
		ctx,
		repo.db,
		`
		INSERT INTO attachments_articles(attachment_id, article_id)
		SELECT $1::bigint AS attachment_id, $2::bigint AS article_id
		FROM articles a
		LEFT JOIN attachments aa ON aa.id = $1::bigint
		WHERE aa.user_id = $3::bigint
			AND aa.deleted_at IS NULL
		  AND a.id = $2::bigint
		RETURNING *;
		`,
		attachment_id,
		article_id,
		user_id,
	)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return ErrBindAttachment
	}
	return nil
}

func (repo *PostgresRepository) DeleteAttachment(ctx context.Context, attachment_id, user_id int64) error {
	tx := repo.db.MustBegin()
	_, err := pu.DispatchTx[dto.AttachmentArticleDbo](
		ctx,
		tx,
		`
		DELETE FROM attachments_articles
		USING attachments_articles AS a
		LEFT JOIN attachments aa ON a.attachment_id = aa.id
		LEFT JOIN user_roles ur ON aa.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE attachments_articles.attachment_id = $1::bigint
		  AND (aa.user_id = $2::bigint OR r.name = 'admin');
		`,
		attachment_id,
		user_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = pu.DispatchTx[dto.AttachmentDbo](
		ctx,
		tx,
		`
		DELETE FROM attachments
		USING attachments AS a
		LEFT JOIN user_roles ur ON a.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE attachments.id = $1::bigint
		  AND (a.user_id = $2::bigint OR r.name = 'admin');
		`,
		attachment_id,
		user_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

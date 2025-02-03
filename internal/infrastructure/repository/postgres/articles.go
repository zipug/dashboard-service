package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrGetArticleById = errors.New("could not get article by id")
	ErrGetAllArticles = errors.New("could not get all articles")
	ErrCreateArticle  = errors.New("could not create article")
	ErrUpdateArticle  = errors.New("could not update article")
	ErrDeleteArticle  = errors.New("could not delete article")
)

func (repo *PostgresRepository) GetArticleById(ctx context.Context, article_id, user_id int64) (*dto.ArticleDbo, error) {
	rows, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		SELECT a.id,
					 a.name,
					 a.description,
					 a.content,
					 a.project_id
		FROM articles a
		LEFT JOIN projects p ON a.project_id = p.id
		LEFT JOIN user_roles ur ON p.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE a.id = $1::bigint
		  AND a.deleted_at IS NULL
		  AND (p.user_id = $2::bigint OR r.name = 'admin');
		`,
		article_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGetArticleById
	}
	row := rows[0]
	return &row, nil
}

func (repo *PostgresRepository) GetAllArticles(ctx context.Context, user_id int64) ([]dto.ArticleDbo, error) {
	rows, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		SELECT a.id,
					 a.name,
					 a.description,
					 a.content,
					 a.project_id
		FROM articles a
		LEFT JOIN projects p ON a.project_id = p.id
		LEFT JOIN user_roles ur ON p.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE a.deleted_at IS NULL
		  AND (p.user_id = $1::bigint OR r.name = 'admin');
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrGetAllArticles
	}
	return rows, nil
}

func (repo *PostgresRepository) CreateArticle(ctx context.Context, article dto.ArticleDbo) (int64, error) {
	rows, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		INSERT INTO articles (name, description, content, project_id)
		VALUES ($1::text, $2::text, $3::text, $4::bigint)
		RETURNING *;
		`,
		article.Name,
		article.Description,
		article.Content,
		article.ProjectId,
	)
	if err != nil {
		return -1, err
	}
	if len(rows) == 0 {
		return -1, ErrCreateArticle
	}
	return rows[0].Id, nil
}

func (repo *PostgresRepository) UpdateArticle(ctx context.Context, article dto.ArticleDbo, user_id int64) (*dto.ArticleDbo, error) {
	rows, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		UPDATE articles a
		SET name = COALESCE(NULLIF($1::text, ''), t.name),
				description = COALESCE(NULLIF($2::text, ''), t.description),
				content = COALESCE(NULLIF($3::text, ''), t.content),
				project_id = COALESCE(NULLIF($4::bigint, 0), t.project_id)
		FROM (
			SELECT at.id, at.name, at.description, at.content, at.project_id
			FROM articles at
			LEFT JOIN projects p ON at.project_id = p.id
			LEFT JOIN user_roles ur ON p.user_id = ur.user_id
			LEFT JOIN roles r ON ur.role_id = r.id
			WHERE at.id = $5::bigint
		    AND (p.user_id = $6::bigint OR r.name = 'admin')
		) AS t(id, name, description, content, project_id)
		WHERE t.id = a.id
		RETURNING a.*;
		`,
		article.Name,
		article.Description,
		article.Content,
		article.ProjectId,
		article.Id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrUpdateArticle
	}
	art := rows[0]
	return &art, nil
}

func (repo *PostgresRepository) DeleteArticle(ctx context.Context, article_id, user_id int64) error {
	_, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		DELETE FROM articles
		USING articles AS a
		LEFT JOIN projects p ON a.project_id = p.id
		LEFT JOIN user_roles ur ON p.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE articles.id = $1::bigint
		  AND (p.user_id = $2::bigint OR r.name = 'admin');
		`,
		article_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

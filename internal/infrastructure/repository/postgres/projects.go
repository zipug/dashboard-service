package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	"database/sql"
	"encoding/json"
	"errors"
)

var (
	ErrGetProjectById = errors.New("could not get project by id")
	ErrGetProjects    = errors.New("could not get projects")
	ErrCreateProject  = errors.New("could not create project")
	ErrUpdateProject  = errors.New("could not update project")
	ErrDeleteProject  = errors.New("could not delete project")
)

type GetProjectByIdParams struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	AvatarUrl   sql.NullString `db:"avatar_url" json:"avatar_url"`
	RemoteUrl   string         `db:"remote_url" json:"remote_url"`
	UserId      int64          `db:"user_id" json:"user_id"`
	Content     string         `db:"content" json:"content"`
}

func (repo *PostgresRepository) GetProjectById(ctx context.Context, project_id int64) (*dto.ProjectsContentDbo, error) {
	rows, err := repo.db.QueryxContext(
		ctx,
		`
		SELECT p.id,
					 p.name,
					 p.description,
					 p.avatar_url,
					 p.remote_url,
					 p.user_id,
					 jsonb_build_object(
						 'id', a.id,
						 'name', a.name,
						 'description', a.description,
						 'article_url', a.article_url,
						 'project_id', a.project_id
					 )::text AS content
		FROM projects p
		LEFT JOIN articles a ON p.id = a.project_id
		WHERE p.id = $1::bigint
		  AND p.deleted_at IS NULL
		  AND a.deleted_at IS NULL;
		`,
		project_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var proj dto.ProjectDbo
	var content []dto.ArticleDbo
	for rows.Next() {
		var content_el dto.ArticleNullSafeDbo
		c := GetProjectByIdParams{}
		if err := rows.StructScan(&c); err != nil {
			return nil, err
		}
		proj = dto.ProjectDbo{
			Id:          c.Id,
			Name:        c.Name,
			Description: c.Description,
			AvatarUrl:   c.AvatarUrl,
			RemoteUrl:   c.RemoteUrl,
			UserId:      c.UserId,
		}
		if err := json.Unmarshal([]byte(c.Content), &content_el); err != nil {
			return nil, err
		}
		if content_el.Id != 0 {
			content = append(content, content_el.ToDbo())
		} else {
			return nil, errors.New("could not append article to project")
		}
	}
	res := dto.ProjectsContentDbo{
		Project: proj,
		Content: content,
	}
	return &res, nil
}

func (repo *PostgresRepository) GetAllProjects(ctx context.Context) ([]dto.ProjectsContentDbo, error) {
	return []dto.ProjectsContentDbo{}, nil
}

func (repo *PostgresRepository) CreateProject(ctx context.Context, project dto.ProjectDbo) (int64, error) {
	return -1, nil
}

func (repo *PostgresRepository) UpdateProject(ctx context.Context, project dto.ProjectDbo) (*dto.ProjectDbo, error) {
	return nil, nil
}

func (repo *PostgresRepository) DeleteProject(ctx context.Context, project_id int64) error {
	return nil
}

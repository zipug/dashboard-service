package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
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

func (repo *PostgresRepository) GetProjectById(ctx context.Context, project_id, user_id int64) (*dto.ProjectsContentDbo, error) {
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
						 'content', a.content,
						 'project_id', a.project_id
					 )::text AS content
		FROM projects p
		LEFT JOIN articles a ON p.id = a.project_id
		LEFT JOIN user_roles ur ON ur.user_id = $2::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		LEFT JOIN users u ON u.id = ur.user_id
		WHERE p.id = $1::bigint
		  AND (p.user_id = $2::bigint OR p.user_id = u.created_by OR r.name = 'admin')
		  AND p.deleted_at IS NULL
		  AND a.deleted_at IS NULL;
		`,
		project_id,
		user_id,
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
		}
	}
	res := dto.ProjectsContentDbo{
		Project: proj,
		Content: content,
	}
	return &res, nil
}

func (repo *PostgresRepository) GetAllProjects(ctx context.Context, user_id int64) ([]dto.ProjectsContentDbo, error) {
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
						 'content', a.content,
						 'project_id', a.project_id
					 )::text AS content
		FROM projects p
		LEFT JOIN articles a ON p.id = a.project_id
		LEFT JOIN user_roles ur ON ur.user_id = $1::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		LEFT JOIN users u ON u.id = ur.user_id
		WHERE p.deleted_at IS NULL
		  AND a.deleted_at IS NULL
		  AND (p.user_id = $1::bigint OR p.user_id = u.created_by OR r.name = 'admin');
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []dto.ProjectsContentDbo
	content_map := make(map[int64]struct {
		Project dto.ProjectDbo
		Content []dto.ArticleDbo
	})
	for rows.Next() {
		var content_el dto.ArticleNullSafeDbo
		c := GetProjectByIdParams{}
		if err := rows.StructScan(&c); err != nil {
			return nil, err
		}
		proj := dto.ProjectDbo{
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
			arr := content_map[proj.Id].Content
			content_map[proj.Id] = struct {
				Project dto.ProjectDbo
				Content []dto.ArticleDbo
			}{
				proj,
				append(arr, content_el.ToDbo()),
			}
		} else {
			content_map[proj.Id] = struct {
				Project dto.ProjectDbo
				Content []dto.ArticleDbo
			}{
				proj,
				nil,
			}
		}
	}
	for _, val := range content_map {
		res = append(res, dto.ProjectsContentDbo{
			Project: val.Project,
			Content: val.Content,
		})
	}
	if len(res) < 1 {
		return res, ErrGetProjects
	}
	return res, nil
}

func (repo *PostgresRepository) CreateProject(ctx context.Context, project dto.ProjectDbo) (int64, error) {
	project_rows, err := pu.Dispatch[dto.ProjectDbo](
		ctx,
		repo.db,
		`
		INSERT INTO projects (name, description, avatar_url, remote_url, user_id)
		VALUES ($1::text, $2::text, $3::text, $4::text, $5::bigint)
		RETURNING *;
		`,
		project.Name,
		project.Description,
		project.AvatarUrl,
		project.RemoteUrl,
		project.UserId,
	)
	if err != nil {
		return -1, err
	}
	if len(project_rows) == 0 {
		return -1, ErrCreateProject
	}
	proj := project_rows[0]
	return proj.Id, nil
}

func (repo *PostgresRepository) UpdateProject(ctx context.Context, project dto.ProjectDbo, user_id int64) (*dto.ProjectDbo, error) {
	project_rows, err := pu.Dispatch[dto.ProjectDbo](
		ctx,
		repo.db,
		`
		UPDATE projects p
		SET name = COALESCE(NULLIF($1::text, ''), t.name),
			description = COALESCE(NULLIF($2::text, ''), t.description),
			avatar_url = COALESCE(NULLIF($3::text, ''), t.avatar_url),
			remote_url = COALESCE(NULLIF($4::text, ''), t.remote_url)
		FROM (
			SELECT pt.id, pt.name, pt.description, pt.avatar_url, pt.remote_url
			FROM projects pt
			WHERE pt.id = $5::bigint
		    AND pt.user_id = $6::bigint
		) AS t(id, name, description, avatar_url, remote_url)
		WHERE t.id = p.id
		RETURNING *;
		`,
		project.Name,
		project.Description,
		project.AvatarUrl,
		project.RemoteUrl,
		project.Id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(project_rows) == 0 {
		return nil, ErrUpdateProject
	}
	proj := project_rows[0]
	return &proj, nil
}

func (repo *PostgresRepository) DeleteProject(ctx context.Context, project_id, user_id int64) error {
	_, err := pu.Dispatch[dto.ProjectDbo](
		ctx,
		repo.db,
		`
		DELETE FROM projects
		USING projects AS p
		LEFT JOIN user_roles ur ON ur.user_id = $2::bigint
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE projects.id = $1::bigint
		  AND (p.user_id = $2::bigint OR r.name = 'admin');
		`,
		project_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

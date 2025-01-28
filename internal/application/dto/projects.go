package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type ProjectDto struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	RemoteUrl   string `json:"remote_url,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
}

type ProjectDbo struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	AvatarUrl   sql.NullString `db:"avatar_url"`
	RemoteUrl   string         `db:"remote_url"`
	UserId      int64          `db:"user_id"`
	CreatedAt   sql.NullTime   `db:"created_at,omitempty"`
	UpdateAt    sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt    sql.NullTime   `db:"deleted_at,omitempty"`
}

type ProjectNullSafeDbo struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	AvatarUrl   string `db:"avatar_url" json:"avatar_url"`
	RemoteUrl   string `db:"remote_url" json:"remote_url"`
	UserId      int64  `db:"user_id" json:"user_id"`
	CreatedAt   string `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdateAt    string `db:"updated_at,omitempty" json:"update_at,omitempty"`
	DeleteAt    string `db:"deleted_at,omitempty" json:"delete_at,omitempty"`
}

func (p *ProjectNullSafeDbo) ToDbo() ProjectDbo {
	return ProjectDbo{
		Id:          p.Id,
		Name:        p.Name,
		Description: sql.NullString{String: p.Description, Valid: true},
		AvatarUrl:   sql.NullString{String: p.AvatarUrl, Valid: true},
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func (p *ProjectDto) ToDbo() ProjectDbo {
	return ProjectDbo{
		Id:          p.Id,
		Name:        p.Name,
		Description: sql.NullString{String: p.Description, Valid: true},
		AvatarUrl:   sql.NullString{String: p.AvatarUrl, Valid: true},
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func (p *ProjectDto) ToValue() models.Project {
	return models.Project{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		AvatarUrl:   p.AvatarUrl,
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func (p *ProjectDbo) ToDto() ProjectDto {
	return ProjectDto{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description.String,
		AvatarUrl:   p.AvatarUrl.String,
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func (p *ProjectDbo) ToValue() models.Project {
	return models.Project{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description.String,
		AvatarUrl:   p.AvatarUrl.String,
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

type ProjectsContentDto struct {
	Project ProjectDto   `json:"project"`
	Content []ArticleDto `json:"content"`
}

func (p *ProjectsContentDto) ToDbo() ProjectsContentDbo {
	var content []ArticleDbo
	for _, article := range p.Content {
		content = append(content, article.ToDbo())
	}
	return ProjectsContentDbo{
		Project: p.Project.ToDbo(),
		Content: content,
	}
}

func (p *ProjectsContentDto) ToValue() models.ProjectsContent {
	var content []models.Article
	for _, article := range p.Content {
		content = append(content, article.ToValue())
	}
	return models.ProjectsContent{
		Project: p.Project.ToValue(),
		Content: content,
	}
}

type ProjectsContentDbo struct {
	Project ProjectDbo   `db:"project"`
	Content []ArticleDbo `db:"content"`
}

func (p *ProjectsContentDbo) ToDbo() ProjectsContentDto {
	var content []ArticleDto
	for _, article := range p.Content {
		content = append(content, article.ToDto())
	}
	return ProjectsContentDto{
		Project: p.Project.ToDto(),
		Content: content,
	}
}

func (p *ProjectsContentDbo) ToValue() models.ProjectsContent {
	var content []models.Article
	for _, article := range p.Content {
		content = append(content, article.ToValue())
	}
	return models.ProjectsContent{
		Project: p.Project.ToValue(),
		Content: content,
	}
}

func ToProjectDto(p models.Project) ProjectDto {
	return ProjectDto{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		AvatarUrl:   p.AvatarUrl,
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func ToProjectDbo(p models.Project) ProjectDbo {
	return ProjectDbo{
		Id:          p.Id,
		Name:        p.Name,
		Description: sql.NullString{String: p.Description, Valid: true},
		AvatarUrl:   sql.NullString{String: p.AvatarUrl, Valid: true},
		RemoteUrl:   p.RemoteUrl,
		UserId:      p.UserId,
	}
}

func ToProjectContentDto(p models.ProjectsContent) ProjectsContentDto {
	content := make([]ArticleDto, 0)
	for _, article := range p.Content {
		content = append(content, ToArticleDto(article))
	}
	return ProjectsContentDto{
		Project: ToProjectDto(p.Project),
		Content: content,
	}
}

func ToProjectContentDbo(p models.ProjectsContent) ProjectsContentDbo {
	content := make([]ArticleDbo, 0)
	for _, article := range p.Content {
		content = append(content, ToArticleDbo(article))
	}
	return ProjectsContentDbo{
		Project: ToProjectDbo(p.Project),
		Content: content,
	}
}

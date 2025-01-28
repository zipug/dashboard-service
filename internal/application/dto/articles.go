package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type ArticleDto struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ArticleUrl  string `json:"article_url,omitempty"`
	ProjectId   int64  `json:"project_id,omitempty"`
}

type ArticleDbo struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	ArticleUrl  string         `db:"article_url"`
	ProjectId   int64          `db:"project_id"`
	CreatedAt   sql.NullTime   `db:"created_at,omitempty"`
	UpdateAt    sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt    sql.NullTime   `db:"deleted_at,omitempty"`
}

type ArticleNullSafeDbo struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ArticleUrl  string `db:"article_url" json:"article_url"`
	ProjectId   int64  `db:"project_id" json:"project_id"`
	CreatedAt   string `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdateAt    string `db:"updated_at,omitempty" json:"update_at,omitempty"`
	DeleteAt    string `db:"deleted_at,omitempty" json:"delete_at,omitempty"`
}

func (a *ArticleNullSafeDbo) ToDbo() ArticleDbo {
	return ArticleDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func (a *ArticleDto) ToDbo() ArticleDbo {
	return ArticleDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func (a *ArticleDbo) ToDto() ArticleDto {
	return ArticleDto{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func (a *ArticleDto) ToValue() models.Article {
	return models.Article{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func (a *ArticleDbo) ToValue() models.Article {
	return models.Article{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func ToArticleDto(a models.Article) ArticleDto {
	return ArticleDto{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

func ToArticleDbo(a models.Article) ArticleDbo {
	return ArticleDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		ArticleUrl:  a.ArticleUrl,
		ProjectId:   a.ProjectId,
	}
}

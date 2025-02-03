package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type ArticlesService interface {
	GetArticleById(context.Context, int64, int64) (models.Article, error)
	GetAllArticles(context.Context, int64) ([]models.Article, error)
	CreateArticle(context.Context, models.Article) (int64, error)
	UpdateArticle(context.Context, models.Article, int64) (models.Article, error)
	DeleteArticle(context.Context, int64, int64) error
}

type ArticleRepository interface {
	GetArticleById(context.Context, int64, int64) (*dto.ArticleDbo, error)
	GetAllArticles(context.Context, int64) ([]dto.ArticleDbo, error)
	CreateArticle(context.Context, dto.ArticleDbo) (int64, error)
	UpdateArticle(context.Context, dto.ArticleDbo, int64) (*dto.ArticleDbo, error)
	DeleteArticle(context.Context, int64, int64) error
}

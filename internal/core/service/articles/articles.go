package articles

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type ArticlesService struct {
	repo ports.ArticleRepository
}

func NewArticlesService(repo ports.ArticleRepository) *ArticlesService {
	return &ArticlesService{repo: repo}
}

func (r *ArticlesService) GetArticleById(ctx context.Context, article_id, user_id int64) (models.Article, error) {
	articleDbo, err := r.repo.GetArticleById(ctx, article_id, user_id)
	if err != nil {
		return models.Article{}, err
	}
	article := articleDbo.ToValue()
	return article, nil
}

func (r *ArticlesService) GetAllArticles(ctx context.Context, user_id int64) ([]models.Article, error) {
	articlesDbo, err := r.repo.GetAllArticles(ctx, user_id)
	if err != nil {
		return []models.Article{}, err
	}
	var articles []models.Article
	for _, article := range articlesDbo {
		articles = append(articles, article.ToValue())
	}
	return articles, nil
}

func (r *ArticlesService) CreateArticle(ctx context.Context, article models.Article) (int64, error) {
	articleDbo := dto.ToArticleDbo(article)
	article_id, err := r.repo.CreateArticle(ctx, articleDbo)
	if err != nil {
		return article_id, err
	}
	return article_id, nil
}

func (r *ArticlesService) UpdateArticle(ctx context.Context, article models.Article, user_id int64) (models.Article, error) {
	articleDbo := dto.ToArticleDbo(article)
	newArticle, err := r.repo.UpdateArticle(ctx, articleDbo, user_id)
	if err != nil {
		return models.Article{}, err
	}
	return newArticle.ToValue(), nil
}

func (r *ArticlesService) DeleteArticle(ctx context.Context, article_id, user_id int64) error {
	return r.repo.DeleteArticle(ctx, article_id, user_id)
}

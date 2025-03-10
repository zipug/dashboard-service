package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"errors"
)

func (d *DashboardService) GetArticleById(article_id, user_id int64) (models.Article, error) {
	ctx := context.Background()
	d.log.Log("info", "getting article by id", logger.WithInt64Attr("article_id", article_id))
	article, err := d.article.GetArticleById(ctx, article_id, user_id)
	if err != nil {
		d.log.Log("error", "error while getting article by id", logger.WithErrAttr(err))
		return article, err
	}
	attachments, err := d.attachment.GetAllAttachmentsByArticleId(ctx, article_id)
	if err != nil {
		d.log.Log("error", "error while getting attachments by article id", logger.WithErrAttr(err))
		return article, err
	}
	var object_ids []string
	for _, attachment := range attachments {
		object_ids = append(object_ids, attachment.ObjectId)
	}
	data, err := d.minio.GetManyFileUrls(ctx, object_ids, "attachments")
	if err != nil {
		d.log.Log("error", "error while getting attachments by article id", logger.WithErrAttr(err))
		return article, err
	}
	for i := 0; i < len(attachments); i++ {
		article.Attachments = append(article.Attachments, models.Attachment{
			Id:          attachments[i].Id,
			Name:        attachments[i].Name,
			Description: attachments[i].Description,
			Mimetype:    attachments[i].Mimetype,
			URL:         data[attachments[i].ObjectId].Url,
		})
	}
	d.log.Log("info", "article successfully get", logger.WithStrAttr("article_id", article.Name))
	return article, nil
}

func (d *DashboardService) GetAllArticles(user_id int64) ([]models.Article, error) {
	ctx := context.Background()
	d.log.Log("info", "getting all articles")
	articles, err := d.article.GetAllArticles(ctx, user_id)
	if err != nil {
		d.log.Log("error", "error while fetching articles", logger.WithErrAttr(err))
		return articles, err
	}
	for i := 0; i < len(articles); i++ {
		attachments, err := d.attachment.GetAllAttachmentsByArticleId(ctx, articles[i].Id)
		if err != nil {
			d.log.Log("error", "error while getting attachments by article id", logger.WithErrAttr(err))
			return nil, err
		}
		var object_ids []string
		for _, attachment := range attachments {
			object_ids = append(object_ids, attachment.ObjectId)
		}
		data, err := d.minio.GetManyFileUrls(ctx, object_ids, "attachments")
		if err != nil {
			d.log.Log("error", "error while getting attachments by article id", logger.WithErrAttr(err))
			return nil, err
		}
		for j := 0; j < len(attachments); j++ {
			articles[i].Attachments = append(articles[i].Attachments, models.Attachment{
				Id:          attachments[j].Id,
				Name:        attachments[j].Name,
				Description: attachments[j].Description,
				Mimetype:    attachments[j].Mimetype,
				URL:         data[attachments[j].ObjectId].Url,
			})
		}
	}
	d.log.Log(
		"info",
		"articles successfully fetched",
		logger.WithInt64Attr("article_count", int64(len(articles))),
	)
	return articles, nil
}

func (d *DashboardService) CreateArticle(article models.Article) (int64, error) {
	ctx := context.Background()
	d.log.Log("info", "creating article")
	article_id, err := d.article.CreateArticle(ctx, article)
	if err != nil {
		d.log.Log("error", "error while creating article", logger.WithErrAttr(err))
		return article_id, err
	}
	d.log.Log("info", "article successfully created", logger.WithStrAttr("article_id", article.Name))
	return article_id, nil
}

func (d *DashboardService) UpdateArticle(article models.Article, user_id int64) (models.Article, error) {
	ctx := context.Background()
	d.log.Log("info", "updating article", logger.WithInt64Attr("article_id", article.Id))
	new_article, err := d.article.UpdateArticle(ctx, article, user_id)
	if err != nil {
		d.log.Log("error", "error while updating article", logger.WithErrAttr(err))
		return new_article, err
	}
	d.log.Log("info", "article successfully updated", logger.WithStrAttr("article_id", article.Name))
	return new_article, nil
}

func (d *DashboardService) DeleteArticle(article_id, user_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "deleting article", logger.WithInt64Attr("project_id", article_id))
	err := d.article.DeleteArticle(ctx, article_id, user_id)
	if err != nil {
		d.log.Log("error", "error while deleting article", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "article successfully deleted", logger.WithInt64Attr("project_id", article_id))
	return nil
}

func (d *DashboardService) UploadArticles(name, extension, content string) error {
	if extension != "html" {
		return errors.New("invalid file extension")
	}
	return nil
}

package updatearticle

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/render"
)

type DashboardService interface {
	UpdateArticle(models.Article) (models.Article, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func UpdateArticle(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		var req dto.ArticleDto
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		article, err := app.UpdateArticle(req.ToValue())
		if err != nil {
			resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to save project info"}}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: dto.ToArticleDto(article)})
	}
}

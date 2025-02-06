package bindattachment

import (
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/render"
)

type DashboardService interface {
	BindAttachment(int64, int64, int64) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(*http.Request) map[string]interface{}
}

type Dto struct {
	AttachmentId int64 `json:"attachment_id"`
	ArticleId    int64 `json:"article_id"`
}

func BindAttachment(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}

		var req Dto

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}

		if err := app.BindAttachment(req.AttachmentId, req.ArticleId, int64(authUserId)); err != nil {
			log.Log("error", "error while binding attachment", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while binding attachment"}})
			return
		}

		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

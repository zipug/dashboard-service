package resolvechat

import (
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type DashboardService interface {
	ResolveQuestion(chat_id int64) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func ResolveQuestion(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		query_id := chi.URLParam(r, "id")
		if query_id == "" {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting id"}})
			return
		}
		id, err := strconv.ParseInt(query_id, 10, 64)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting id"}})
			return
		}
		if err := app.ResolveQuestion(id); err != nil {
			errs := strings.Split(err.Error(), "\n")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: errs})
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

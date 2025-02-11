package runbotbyid

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
	RunBotById(int64, int64) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(*http.Request) map[string]interface{}
}

func RunBot(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}
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
		if err := app.RunBotById(id, int64(authUserId)); err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

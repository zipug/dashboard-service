package grantroletouser

import (
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	GrantRoleToUser(int64, int64) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func GrantRoleToUser(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.URL.Query().Get("user_id") == "" || r.URL.Query().Get("role_id") == "" {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting id"}})
			return
		}
		query_user_id := r.URL.Query().Get("user_id")
		user_id, err := strconv.ParseInt(query_user_id, 10, 64)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting user_id"}})
			return
		}
		query_role_id := r.URL.Query().Get("role_id")
		role_id, err := strconv.ParseInt(query_role_id, 10, 64)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting role_id"}})
			return
		}
		if err := app.GrantRoleToUser(user_id, role_id); err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

package createproject

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	CreateProject(models.Project) (int64, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(*http.Request) map[string]interface{}
}

func CreateProject(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}
		var req dto.ProjectDto
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		payload := req.ToValue()
		payload.Id = int64(authUserId)
		project_id, err := app.CreateProject(payload)
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: project_id})
	}
}

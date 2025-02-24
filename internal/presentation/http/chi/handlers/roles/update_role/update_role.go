package updaterole

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/render"
)

type DashboardService interface {
	UpdateRole(models.Role) (models.Role, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func UpdateRole(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var req dto.RoleDto
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		role, err := app.UpdateRole(req.ToValue())
		if err != nil {
			resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to save role info"}}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: dto.ToRoleDto(role)})
	}
}

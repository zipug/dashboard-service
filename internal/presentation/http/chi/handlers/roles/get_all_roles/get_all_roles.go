package getallroles

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
	GetAllRoles() ([]models.Role, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func GetAllRoles(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		roles, err := app.GetAllRoles()
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		var resp []dto.RoleDto
		for _, role := range roles {
			resp = append(resp, dto.ToRoleDto(role))
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: resp})
	}
}

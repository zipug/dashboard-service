package updateroleperms

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/render"
)

type DashboardService interface {
	UpdateRolePerms(int64, []models.PermissionData) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type UpdateRolePermsPayload struct {
	Perms  []dto.PermissionDto `json:"perms"`
	RoleId int64               `json:"role_id"`
}

func UpdateRolePerms(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var req UpdateRolePermsPayload
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		var perms []models.PermissionData
		for _, perm := range req.Perms {
			perms = append(perms, perm.ToValue())
		}
		err := app.UpdateRolePerms(req.RoleId, perms)
		if err != nil {
			resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to update permissions"}}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

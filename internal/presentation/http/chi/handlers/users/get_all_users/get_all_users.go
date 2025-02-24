package getallusers

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
	GetAllUsers() ([]models.User, error)
	GetRoleByUserId(user_id int64) (models.Role, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func GetAllUsers(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		users, err := app.GetAllUsers()
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		var resp []dto.UserDto
		for _, user := range users {
			role, err := app.GetRoleByUserId(int64(user.Id))
			if err != nil {
				errs := strings.Split(err.Error(), "\n")
				resp := handlers.Response{Status: handlers.Failed, Errors: errs}
				render.JSON(w, r, resp)
				return
			}
			dtoUser := dto.ToUserDto(user)
			dtoUser.Role = dto.ToRoleDto(role)
			resp = append(resp, dtoUser)
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: resp})
	}
}

package getuserinfo

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
	GetUserById(id int64) (models.User, error)
	GetRoleByUserId(user_id int64) (models.Role, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(r *http.Request) map[string]interface{}
}

func GetUserInfo(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}
		user, err := app.GetUserById(int64(authUserId))
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		role, err := app.GetRoleByUserId(int64(authUserId))
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		userDto := dto.ToSafeUserDto(user)
		userDto.Role = dto.ToRoleDto(role)
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: userDto})
	}
}

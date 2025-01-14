package login

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type DashboardService interface {
	LoginUser(dto.UserDto) (models.User, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetTokenAuth() *jwtauth.JWTAuth
	HashPassword(string) (string, error)
	ValidatePassword(string, string) bool
}

type RegisterResponse struct {
	Id int64 `json:"id"`
}

func LoginUser(app DashboardService, log Logger, auth Auth, accessTokenExp time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserDto
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		user, err := app.LoginUser(req)
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		if ok := auth.ValidatePassword(string(user.Password), req.Password); !ok {
			log.Log("error", "invalid user password", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"password does not match"}})
			return
		}
		claims := map[string]interface{}{"user_id": user.Id, "exp": time.Now().Add(accessTokenExp).Unix()}
		_, token, err := auth.GetTokenAuth().Encode(claims)
		if err != nil {
			log.Log("error", "error while generate jwt token", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"something went wrong"}})
			return
		}
		authToken := dto.AuthenticateDto{Token: token}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: authToken})
	}
}

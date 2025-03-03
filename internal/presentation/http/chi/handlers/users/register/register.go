package register

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type DashboardService interface {
	RegisterUser(dto.UserDto) (int64, error)
	RegisterSupport(dto.UserDto, int64) (int64, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(r *http.Request) map[string]interface{}
	GetTokenAuth() *jwtauth.JWTAuth
	HashPassword(string) (string, error)
}

func RegisterUser(app DashboardService, log Logger, auth Auth, accessTokenExp time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var req dto.UserDto
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		req.Password, err = auth.HashPassword(req.Password)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"something went wrong"}})
			return
		}
		req.RepeatPassword, err = auth.HashPassword(req.RepeatPassword)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"something went wrong"}})
			return
		}
		id, err := func() (int64, error) {
			authClaims := auth.GetClaims(r)
			authUserId, ok := authClaims["user_id"].(float64)
			if !ok {
				log.Log("error", "invalid user_id in jwt token")
				return app.RegisterUser(req)
			}
			return app.RegisterSupport(req, int64(authUserId))
		}()
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		claims := map[string]interface{}{"user_id": id, "exp": time.Now().Add(accessTokenExp).Unix()}
		_, token, err := auth.GetTokenAuth().Encode(claims)
		if err != nil {
			log.Log("error", "error while generate jwt token", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"something went wrong"}})
			return
		}

		authToken := dto.AuthenticateDto{Token: token}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: authToken})
	}
}

package verifyuser

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/render"
)

type DashboardService interface {
	VerifyUser(int64, models.OTPCode) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(r *http.Request) map[string]interface{}
}

func Verify(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}

		var req dto.VerifyUserDto
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		if err := app.VerifyUser(int64(authUserId), req.Code); err != nil {
			resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to verify your otp code"}}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: []string{"successfully verified"}})
	}
}

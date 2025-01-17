package sendotp

import (
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type DashboardService interface {
	SendOTPByUserId(int64, models.OTPTarget) error
	SendOTPByUserEmail(string, models.OTPTarget) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(r *http.Request) map[string]interface{}
}

func SendOTP(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		query_target := chi.URLParam(r, "target")
		if query_target == "" {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting targer"}})
			return
		}
		target := models.OTPTarget(query_target)
		switch target {
		case models.VERIFICATION:
			authClaims := auth.GetClaims(r)
			authUserId, ok := authClaims["user_id"].(float64)
			if !ok {
				log.Log("error", "invalid user_id in jwt token")
				render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
				return
			}
			err := app.SendOTPByUserId(int64(authUserId), target)
			if err != nil {
				resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to send otp code"}}
				render.JSON(w, r, resp)
				return
			}
			render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: []string{"Your verification code has been sent"}})
		case models.AUTHENTICATION:
			var req dto.UserDto
			if err := render.DecodeJSON(r.Body, &req); err != nil {
				log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
				render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
				return
			}
			if err := app.SendOTPByUserEmail(req.Email, target); err != nil {
				resp := handlers.Response{Status: handlers.Failed, Errors: []string{"failed to send otp code"}}
				render.JSON(w, r, resp)
				return
			}
			render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: []string{"Your verification code has been sent"}})
		default:
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid otp code sending target"}})
		}
	}
}

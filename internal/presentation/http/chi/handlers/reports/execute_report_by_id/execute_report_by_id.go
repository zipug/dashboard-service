package executereportbyid

import (
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	ExecuteReport(int64, int64, string, string) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(*http.Request) map[string]interface{}
}

type Payload struct {
	ReportID int64  `json:"report_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func ExecuteReport(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}
		var req Payload
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Log("error", "error while decoding request body", logger.WithErrAttr(err))
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while decoding request body"}})
			return
		}
		if err := app.ExecuteReport(req.ReportID, int64(authUserId), req.DateFrom, req.DateTo); err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

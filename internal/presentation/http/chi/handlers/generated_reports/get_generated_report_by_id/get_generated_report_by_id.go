package getgeneratedreportbyid

import (
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/presentation/http/chi/handlers"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type DashboardService interface {
	GetGeneratedReportById(int64, int64) (models.GeneratedReport, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(r *http.Request) map[string]interface{}
}

func GetGeneratedReportById(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/csv; charset=utf-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}
		query_id := chi.URLParam(r, "id")
		if query_id == "" {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting id"}})
			return
		}
		id, err := strconv.ParseInt(query_id, 10, 64)
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting id"}})
			return
		}
		generated_report, err := app.GetGeneratedReportById(id, int64(authUserId))
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", generated_report.ObjectID))

		writter := csv.NewWriter(w)
		defer writter.Flush()

		if err := writter.Write(generated_report.Content.Headers); err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while downloading report", err.Error()}})
			return
		}
		for _, rec := range generated_report.Content.Rows {
			if err := writter.Write(rec); err != nil {
				render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while downloading report", err.Error()}})
				return
			}
		}
	}
}

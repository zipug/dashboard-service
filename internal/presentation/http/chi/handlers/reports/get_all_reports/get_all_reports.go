package getallreports

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
	GetAllReports() ([]models.Report, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func GetAllReports(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		reports, err := app.GetAllReports()
		if err != nil {
			errs := strings.Split(err.Error(), "\n")
			resp := handlers.Response{Status: handlers.Failed, Errors: errs}
			render.JSON(w, r, resp)
			return
		}
		var resp []dto.ReportDto
		for _, report := range reports {
			resp = append(resp, dto.ToReportDto(report))
		}
		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: resp})
	}
}

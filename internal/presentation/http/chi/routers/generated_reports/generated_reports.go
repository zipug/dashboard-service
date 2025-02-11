package generatedreports

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	deletegeneratedreportbyid "dashboard/internal/presentation/http/chi/handlers/generated_reports/delete_generated_report_by_id"
	getallgeneratedreports "dashboard/internal/presentation/http/chi/handlers/generated_reports/get_all_generated_reports"
	getgeneratedreportbyid "dashboard/internal/presentation/http/chi/handlers/generated_reports/get_generated_report_by_id"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func GeneratedReportsRouter(r chi.Router) func(
	*a.DashboardService,
	*l.Logger,
	*auth.Auth,
	*config.AppConfig,
) chi.Router {
	return func(
		app *a.DashboardService,
		log *l.Logger,
		auth *auth.Auth,
		cfg *config.AppConfig,
	) chi.Router {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.GetTokenAuth()))
			r.Use(jwtauth.Authenticator(auth.GetTokenAuth()))
			guard := can.NewGuard()
			guard.AddVerifier(app.ValidateUserPermissions)
			r.
				With(guard.Can(auth.GetTokenAuth(), "generatedreports_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getgeneratedreportbyid.GetGeneratedReportById(app, log, auth))
					r.Get("/all", getallgeneratedreports.GetAllGeneratedReports(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "generatedreports_feature:delete")).
				Delete("/{id}", deletegeneratedreportbyid.DeleteGeneratedReportById(app, log, auth))
		})
		return r
	}
}

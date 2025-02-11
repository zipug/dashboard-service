package reports

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	executereportbyid "dashboard/internal/presentation/http/chi/handlers/reports/execute_report_by_id"
	getallreports "dashboard/internal/presentation/http/chi/handlers/reports/get_all_reports"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func ReportsRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "reports_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/all", getallreports.GetAllReports(app, log))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "reports_feature:update")).
				Post("/execute", executereportbyid.ExecuteReport(app, log, auth))
		})
		return r
	}
}

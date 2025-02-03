package projects

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	createproject "dashboard/internal/presentation/http/chi/handlers/projects/create_project"
	deleteprojectbyid "dashboard/internal/presentation/http/chi/handlers/projects/delete_project_by_id"
	getallprojects "dashboard/internal/presentation/http/chi/handlers/projects/get_all_projects"
	getprojectbyid "dashboard/internal/presentation/http/chi/handlers/projects/get_project_by_id"
	updateproject "dashboard/internal/presentation/http/chi/handlers/projects/update_project"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func ProjectsRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "projects_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getprojectbyid.GetProjectById(app, log, auth))
					r.Get("/all", getallprojects.GetAllProjects(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "projects_feature:create")).
				Post("/create", createproject.CreateProject(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "projects_feature:update")).
				Post("/update", updateproject.UpdateProject(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "projects_feature:delete")).
				Delete("/{id}", deleteprojectbyid.DeleteProject(app, log, auth))
		})
		return r
	}
}

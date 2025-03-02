package roles

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	createrole "dashboard/internal/presentation/http/chi/handlers/roles/create_role"
	deleterolebyid "dashboard/internal/presentation/http/chi/handlers/roles/delete_role_by_id"
	getallroles "dashboard/internal/presentation/http/chi/handlers/roles/get_all_roles"
	getrolebyid "dashboard/internal/presentation/http/chi/handlers/roles/get_role_by_id"
	updaterole "dashboard/internal/presentation/http/chi/handlers/roles/update_role"
	updateroleperms "dashboard/internal/presentation/http/chi/handlers/roles/update_role_permissions"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func RolesRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "roles_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getrolebyid.GetRoleById(app, log))
					r.Get("/all", getallroles.GetAllRoles(app, log))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "roles_feature:create")).
				Post("/create", createrole.CreateRole(app, log))
			r.
				With(guard.Can(auth.GetTokenAuth(), "roles_feature:update")).
				Group(func(r chi.Router) {
					r.Post("/update", updaterole.UpdateRole(app, log))
					r.Patch("/update-permissions", updateroleperms.UpdateRolePerms(app, log))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "roles_feature:delete")).
				Delete("/{id}", deleterolebyid.DeleteRole(app, log))
		})
		return r
	}
}

package users

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	deleteuserbyid "dashboard/internal/presentation/http/chi/handlers/users/delete_user_by_id"
	getallusers "dashboard/internal/presentation/http/chi/handlers/users/get_all_users"
	getuserbyid "dashboard/internal/presentation/http/chi/handlers/users/get_user_by_id"
	getuserinfo "dashboard/internal/presentation/http/chi/handlers/users/get_user_info"
	grantroletouser "dashboard/internal/presentation/http/chi/handlers/users/grant_role_to_user"
	"dashboard/internal/presentation/http/chi/handlers/users/login"
	"dashboard/internal/presentation/http/chi/handlers/users/register"
	saveuserinfo "dashboard/internal/presentation/http/chi/handlers/users/save_user_info"
	sendotp "dashboard/internal/presentation/http/chi/handlers/users/send_otp"
	verifyuser "dashboard/internal/presentation/http/chi/handlers/users/verify_user"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func UsersRouter(r chi.Router) func(
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
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(auth.GetTokenAuth()))
				r.Use(jwtauth.Authenticator(auth.GetTokenAuth()))

				r.Group(func(r chi.Router) {
					guard := can.NewGuard()
					guard.AddVerifier(app.ValidateUserPermissions)
					r.
						With(guard.Can(auth.GetTokenAuth(), "users_feature:read")).
						Get("/all", getallusers.GetAllUsers(app, log))
					r.
						With(guard.Can(auth.GetTokenAuth(), "users_feature:update")).
						Patch("/grant", grantroletouser.GrantRoleToUser(app, log))
					r.
						With(guard.Can(auth.GetTokenAuth(), "users_feature:delete")).
						Delete("/{id}", deleteuserbyid.DeleteUserById(app, log))
				})

				r.Get("/{id}", getuserbyid.GetUserById(app, log))
				r.Get("/me", getuserinfo.GetUserInfo(app, log, auth))
				r.Post("/verify", verifyuser.Verify(app, log, auth))
				r.Post("/update", saveuserinfo.SaveUser(app, log, auth))
			})

			r.Group(func(r chi.Router) {
				r.Post("/register", register.RegisterUser(app, log, auth, cfg.AccessTokenExpiration))
				r.Post("/login", login.LoginUser(app, log, auth, cfg.AccessTokenExpiration))
				r.Post("/send-code", sendotp.SendOTP(app, log))
			})
		})
		return r
	}
}

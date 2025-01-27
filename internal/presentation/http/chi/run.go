package chi

import (
	"dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	logger "dashboard/internal/common/service/logger/zerolog"
	createrole "dashboard/internal/presentation/http/chi/handlers/roles/create_role"
	deleterolebyid "dashboard/internal/presentation/http/chi/handlers/roles/delete_role_by_id"
	getallroles "dashboard/internal/presentation/http/chi/handlers/roles/get_all_roles"
	getrolebyid "dashboard/internal/presentation/http/chi/handlers/roles/get_role_by_id"
	updaterole "dashboard/internal/presentation/http/chi/handlers/roles/update_role"
	updateroleperms "dashboard/internal/presentation/http/chi/handlers/roles/update_role_permissions"
	"dashboard/internal/presentation/http/chi/routers/users"
	"dashboard/pkg/middlewares/can"
	logger_middleware "dashboard/pkg/middlewares/logger"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type HttpServer struct {
	mux *chi.Mux
	cfg *config.AppConfig
}

var tokenAuth *jwtauth.JWTAuth

func NewHttpServer(app *application.DashboardService) *HttpServer {
	router := chi.NewRouter()

	var cfg *config.AppConfig
	var log *logger.Logger
	var auth *auth.Auth
	retries := 5

	app.Run()
	for retries > 0 {
		switch app.GetState() {
		case application.Ready:
			cfg = app.GetConfig()
			log = app.GetLogger()
			auth = app.GetAuth()
			log.Log("info",
				fmt.Sprintf("Starting application on %s:%d%s",
					cfg.Server.Host,
					cfg.Server.Port,
					cfg.Server.DefaultApi,
				),
			)
			retries = 0
		case application.Down:
			log.Log("error", "application failed to start")
			retries--
		default:
		}
	}

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.URLFormat,
		logger_middleware.New(log.GetLogger()),
	)
	api := cfg.Server.DefaultApi

	router.Route(api, func(r chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		r.Route("/users", func(r chi.Router) {
			users.UsersRouter(r)(app, log, auth, cfg)
		})

		r.Route("/roles", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(auth.GetTokenAuth()))
				r.Use(jwtauth.Authenticator(auth.GetTokenAuth()))
				r.Group(func(r chi.Router) {
					guard := can.NewGuard()
					guard.AddVerifier(app.ValidateUserPermissions)
					r.Use(guard.Can(auth.GetTokenAuth(), "roles_feature:read"))
					r.Get("/{id}", getrolebyid.GetRoleById(app, log))
					r.Get("/all", getallroles.GetAllRoles(app, log))
				})
				r.Group(func(r chi.Router) {
					guard := can.NewGuard()
					guard.AddVerifier(app.ValidateUserPermissions)
					r.Use(guard.Can(auth.GetTokenAuth(), "roles_feature:create"))
					r.Post("/create", createrole.CreateRole(app, log))
				})
				r.Group(func(r chi.Router) {
					guard := can.NewGuard()
					guard.AddVerifier(app.ValidateUserPermissions)
					r.Use(guard.Can(auth.GetTokenAuth(), "roles_feature:update"))
					r.Post("/update", updaterole.UpdateRole(app, log))
					r.Patch("/update-permissions", updateroleperms.UpdateRolePerms(app, log))
				})
				r.Group(func(r chi.Router) {
					guard := can.NewGuard()
					guard.AddVerifier(app.ValidateUserPermissions)
					r.Use(guard.Can(auth.GetTokenAuth(), "roles_feature:delete"))
					r.Delete("/{id}", deleterolebyid.DeleteRole(app, log))
				})
			})
			r.Group(func(r chi.Router) {
			})
		})
	})

	return &HttpServer{
		mux: router,
		cfg: cfg,
	}
}

func (s *HttpServer) Start() {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler:      s.mux,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}
	server.ListenAndServe()
}

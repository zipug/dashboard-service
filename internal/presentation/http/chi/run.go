package chi

import (
	"dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/routers/articles"
	"dashboard/internal/presentation/http/chi/routers/attachments"
	"dashboard/internal/presentation/http/chi/routers/bots"
	generatedreports "dashboard/internal/presentation/http/chi/routers/generated_reports"
	"dashboard/internal/presentation/http/chi/routers/projects"
	"dashboard/internal/presentation/http/chi/routers/reports"
	"dashboard/internal/presentation/http/chi/routers/roles"
	"dashboard/internal/presentation/http/chi/routers/users"
	logger_middleware "dashboard/pkg/middlewares/logger"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{cfg.Server.FrontEndUrl},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
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
			roles.RolesRouter(r)(app, log, auth, cfg)
		})
		r.Route("/projects", func(r chi.Router) {
			projects.ProjectsRouter(r)(app, log, auth, cfg)
		})
		r.Route("/articles", func(r chi.Router) {
			articles.ArticlesRouter(r)(app, log, auth, cfg)
		})
		r.Route("/attachments", func(r chi.Router) {
			attachments.AttachmentsRouter(r)(app, log, auth, cfg)
		})
		r.Route("/reports", func(r chi.Router) {
			reports.ReportsRouter(r)(app, log, auth, cfg)
		})
		r.Route("/generated-reports", func(r chi.Router) {
			generatedreports.GeneratedReportsRouter(r)(app, log, auth, cfg)
		})
		r.Route("/bots", func(r chi.Router) {
			bots.BotsRouter(r)(app, log, auth, cfg)
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

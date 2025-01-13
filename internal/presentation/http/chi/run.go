package chi

import (
	"dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	logger "dashboard/internal/common/service/logger/zerolog"
	getuserbyid "dashboard/internal/presentation/http/chi/handlers/get_user_by_id"
	"dashboard/internal/presentation/http/chi/handlers/login"
	"dashboard/internal/presentation/http/chi/handlers/register"
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
		logger_middleware.New(),
	)
	api := cfg.Server.DefaultApi

	router.Route(api, func(r chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(auth.GetTokenAuth()))
				r.Use(jwtauth.Authenticator(auth.GetTokenAuth()))

				r.Get("/{id}", getuserbyid.GetUserById(app, log))
			})

			r.Group(func(r chi.Router) {
				r.Post("/register", register.RegisterUser(app, log, auth, cfg.AccessTokenExpiration))
				r.Post("/login", login.LoginUser(app, log, auth, cfg.AccessTokenExpiration))
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

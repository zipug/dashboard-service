package bots

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	createbot "dashboard/internal/presentation/http/chi/handlers/bots/create_bot"
	deletebotbyid "dashboard/internal/presentation/http/chi/handlers/bots/delete_bot_by_id"
	getallbots "dashboard/internal/presentation/http/chi/handlers/bots/get_all_bots"
	getbotbyid "dashboard/internal/presentation/http/chi/handlers/bots/get_bot_by_id"
	setbotstatebyid "dashboard/internal/presentation/http/chi/handlers/bots/set_bot_state_by_id"
	updatebotbyid "dashboard/internal/presentation/http/chi/handlers/bots/update_bot_by_id"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func BotsRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "bots_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getbotbyid.GetBotById(app, log, auth))
					r.Get("/all", getallbots.GetAllBots(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "bots_feature:create")).
				Post("/create", createbot.CreateBot(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "bots_feature:update")).
				Group(func(r chi.Router) {
					r.Post("/update", updatebotbyid.UpdateBotById(app, log, auth))
					r.Post("/set-state", setbotstatebyid.SetBotState(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "bots_feature:delete")).
				Delete("/{id}", deletebotbyid.DeleteBotById(app, log, auth))
		})
		return r
	}
}

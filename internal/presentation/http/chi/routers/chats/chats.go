package chats

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	getallchats "dashboard/internal/presentation/http/chi/handlers/chats/get_all_chats"
	getchatsbyid "dashboard/internal/presentation/http/chi/handlers/chats/get_chats_by_id"
	telegramgetuserbyid "dashboard/internal/presentation/http/chi/handlers/chats/telegram_get_user_by_id"
	telegramsendmessage "dashboard/internal/presentation/http/chi/handlers/chats/telegram_send_message"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func ChatsRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "chats_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getchatsbyid.GetChatById(app, log, auth))
					r.Get("/all", getallchats.GetAllChats(app, log, auth))
					r.Post("/tg-user", telegramgetuserbyid.GetTelegramUserById(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "chats_feature:update")).
				Group(func(r chi.Router) {
					r.Post("/send-message", telegramsendmessage.SendMessage(app, log, auth))
				})
		})
		return r
	}
}

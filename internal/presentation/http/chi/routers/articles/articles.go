package articles

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	createarticle "dashboard/internal/presentation/http/chi/handlers/articles/create_article"
	deletearticlebyid "dashboard/internal/presentation/http/chi/handlers/articles/delete_article_by_id"
	getallarticles "dashboard/internal/presentation/http/chi/handlers/articles/get_all_articles"
	getarticlebyid "dashboard/internal/presentation/http/chi/handlers/articles/get_article_by_id"
	updatearticle "dashboard/internal/presentation/http/chi/handlers/articles/update_article"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func ArticlesRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "articles_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getarticlebyid.GetArticleById(app, log, auth))
					r.Get("/all", getallarticles.GetAllArticles(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "articles_feature:create")).
				Group(func(r chi.Router) {
					r.Post("/create", createarticle.CreateArticle(app, log))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "articles_feature:update")).
				Post("/update", updatearticle.UpdateArticle(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "articles_feature:delete")).
				Delete("/{id}", deletearticlebyid.DeleteArticle(app, log, auth))
		})
		return r
	}
}

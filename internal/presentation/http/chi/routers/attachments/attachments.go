package attachments

import (
	a "dashboard/internal/application"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	l "dashboard/internal/common/service/logger/zerolog"
	bindattachment "dashboard/internal/presentation/http/chi/handlers/attachments/bind_attachment"
	getallattachments "dashboard/internal/presentation/http/chi/handlers/attachments/get_all_attachments"
	getattachmentbyid "dashboard/internal/presentation/http/chi/handlers/attachments/get_attachment_by_id"
	uploadfile "dashboard/internal/presentation/http/chi/handlers/attachments/upload_file"
	deleterolebyid "dashboard/internal/presentation/http/chi/handlers/roles/delete_role_by_id"
	"dashboard/pkg/middlewares/can"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func AttachmentsRouter(r chi.Router) func(
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
				With(guard.Can(auth.GetTokenAuth(), "attachments_feature:read")).
				Group(func(r chi.Router) {
					r.Get("/{id}", getattachmentbyid.GetAttachmentById(app, log, auth))
					r.Get("/all", getallattachments.GetAllAttachments(app, log, auth))
				})
			r.
				With(guard.Can(auth.GetTokenAuth(), "attachments_feature:create")).
				Post("/upload", uploadfile.UploadAttachment(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "attachments_feature:update")).
				Post("/bind", bindattachment.BindAttachment(app, log, auth))
			r.
				With(guard.Can(auth.GetTokenAuth(), "attachments_feature:delete")).
				Delete("/{id}", deleterolebyid.DeleteRole(app, log))
		})
		return r
	}
}

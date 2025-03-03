package uploadfile

import (
	"bytes"
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	UploadAttachment(string, string, []byte, int64) (dto.AttachmentDto, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Auth interface {
	GetClaims(*http.Request) map[string]interface{}
}

func UploadAttachment(app DashboardService, log Logger, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		authClaims := auth.GetClaims(r)
		authUserId, ok := authClaims["user_id"].(float64)
		if !ok {
			log.Log("error", "invalid user_id in jwt token")
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"invalid user_id in jwt token"}})
			return
		}

		r.ParseMultipartForm(32 << 20)
		if r.ContentLength > 32<<20 {
			http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
			return
		}
		var res []dto.AttachmentDto
		files := r.MultipartForm.File["file"]
		for _, fileHeader := range files {
			var buf bytes.Buffer
			defer buf.Reset()
			file, err := fileHeader.Open()
			if err != nil {
				render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting file"}})
				log.Log("error", "error while getting file", logger.WithErrAttr(err))
				return
			}
			defer file.Close()
			name := strings.Split(fileHeader.Filename, ".")
			io.Copy(&buf, file)
			content := buf.String()
			log.Log(
				"info",
				"file uploaded",
				logger.WithStrAttr("name", name[0]),
				logger.WithStrAttr("extension", name[1]),
				logger.WithInt64Attr("size", int64(len(content))),
			)
			attachment, err := app.UploadAttachment(
				fileHeader.Filename,
				fmt.Sprintf(".%s", name[1]),
				buf.Bytes(),
				int64(authUserId),
			)
			if err != nil {
				render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while uploading file"}})
				return
			}
			res = append(res, attachment)
		}

		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: res})
	}
}

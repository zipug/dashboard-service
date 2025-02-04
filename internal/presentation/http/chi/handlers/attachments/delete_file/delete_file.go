package deletefile

import (
	"bytes"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	DeleteAttachment(string, string, string) error
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

func DeleteAttachment(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		r.ParseMultipartForm(32 << 20)
		if r.ContentLength > 32<<20 {
			http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
			return
		}
		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while getting file"}})
			log.Log("error", "error while getting file", logger.WithErrAttr(err))
			return
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		io.Copy(&buf, file)
		content := buf.String()
		buf.Reset()
		log.Log(
			"info",
			"file uploaded",
			logger.WithStrAttr("name", name[0]),
			logger.WithStrAttr("extension", name[1]),
			logger.WithInt64Attr("size", int64(len(content))),
		)
		app.DeleteAttachment(name[0], name[1], content)

		render.JSON(w, r, handlers.Response{Status: handlers.Success})
	}
}

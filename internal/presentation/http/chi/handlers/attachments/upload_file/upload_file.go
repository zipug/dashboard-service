package uploadfile

import (
	"bytes"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/presentation/http/chi/handlers"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type DashboardService interface {
	UploadAttachment(string, string, []byte) (string, error)
}

type Logger interface {
	Log(logger.LoggerAction, string, ...logger.LoggerEvent)
}

type Attachment struct {
	URL string `json:"url"`
}

func UploadAttachment(app DashboardService, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		r.ParseMultipartForm(32 << 20)
		if r.ContentLength > 32<<20 {
			http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
			return
		}
		var buf bytes.Buffer
		defer buf.Reset()
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
		log.Log(
			"info",
			"file uploaded",
			logger.WithStrAttr("name", name[0]),
			logger.WithStrAttr("extension", name[1]),
			logger.WithInt64Attr("size", int64(len(content))),
		)
		url, err := app.UploadAttachment(header.Filename, fmt.Sprintf(".%s", name[1]), buf.Bytes())
		if err != nil {
			render.JSON(w, r, handlers.Response{Status: handlers.Failed, Errors: []string{"error while uploading file"}})
			return
		}

		render.JSON(w, r, handlers.Response{Status: handlers.Success, Data: Attachment{URL: url}})
	}
}

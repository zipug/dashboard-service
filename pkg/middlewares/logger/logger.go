package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

		log := zerolog.New(output).
			With().
			Timestamp().
			Logger().
			Level(zerolog.DebugLevel)

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Str("request_id", middleware.GetReqID(r.Context())).
				Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info().
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Str("duration", time.Since(t1).String()).
					Msg("request completed")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

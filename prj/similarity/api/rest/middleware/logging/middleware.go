package logging

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func NewMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			rw := newStatusRW(w)

			next.ServeHTTP(rw, r)

			slog.Info(
				"request processed",
				"request_id", r.Context().Value(model.RequestIDKey),
				"status", rw.Status(),
				"response_size_bytes", rw.Size(),
				"elapsed_time", time.Since(t),
			)
		})
	}
}

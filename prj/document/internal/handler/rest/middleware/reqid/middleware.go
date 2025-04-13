package reqid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func NewMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(model.RequestIDHeader)

			if rid == "" {
				rid = uuid.NewString()
			}

			r = r.WithContext(context.WithValue(r.Context(), model.RequestIDKey, rid))

			next.ServeHTTP(w, r)
		})
	}
}

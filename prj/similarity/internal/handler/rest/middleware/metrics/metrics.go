package metrics

import (
	"net/http"
	"strconv"
	"time"
)

func NewMiddleware(m Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			rw := newStatusRW(w)

			next.ServeHTTP(rw, r)

			m.IncHTTPRequest(
				r.URL.Path,
				strconv.Itoa(rw.Status()),
				rw.Size(),
				time.Since(t).Seconds(),
			)
		})
	}
}

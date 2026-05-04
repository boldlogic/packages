package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// Metrics оборачивает handler и после ответа записывает HTTP-метрики.
//
// Для label route используется шаблон маршрута из chi, если он доступен,
// иначе используется путь запроса.
func (m Middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		route := routePattern(r)
		m.metrics.RecordRequest(r.Method, route, strconv.Itoa(rw.status), time.Since(start))
	})
}

func routePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if pattern := rctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return r.URL.Path
}

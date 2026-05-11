package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

// Recover перехватывает панику в следующем http.Handler и пишет её в лог.
func (m Middleware) Recover(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseRecorder{ResponseWriter: w}
		defer func() {
			if rec := recover(); rec != nil {
				m.logger.Error("возникла паника",
					zap.Any("rec", rec),
					zap.ByteString("stack", debug.Stack()))
				//rw.status = http.StatusInternalServerError
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

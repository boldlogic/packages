package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRecover(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantBody   string
	}{
		{
			name: "паника_до_записи_ответа_возвращает_500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("причина")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "обычный_ответ_проходит_без_изменений",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte("ok"))
			},
			wantStatus: http.StatusCreated,
			wantBody:   "ok",
		},
		{
			name: "паника_после_записи_заголовка_оставляет_исходный_статус",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				panic("причина")
			},
			wantStatus: http.StatusCreated,
		},
	}

	middleware := NewMiddleware(nil, zap.NewNop())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := middleware.Recover(tt.handler)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			require.NotPanics(t, func() {
				wrapped.ServeHTTP(rec, req)
			})

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}

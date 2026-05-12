package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestWriteResp(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantStatus      int
		wantBody        string
		wantContentType bool
	}{
		{
			name: "успешный_json_ответ",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				s := struct {
					D decimal.Decimal `json:"balance"`
				}{
					D: decimal.RequireFromString("11.01"),
				}
				WriteResp(w, http.StatusOK, s)
			},
			wantStatus:      http.StatusOK,
			wantBody:        `{"balance":"11.01"}`,
			wantContentType: true,
		},
		{
			name: "успешный_json_ответ_со_статусом_created",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				s := struct {
					D decimal.Decimal `json:"balance"`
				}{
					D: decimal.RequireFromString("11.01"),
				}
				WriteResp(w, http.StatusCreated, s)
			},
			wantStatus:      http.StatusCreated,
			wantBody:        `{"balance":"11.01"}`,
			wantContentType: true,
		},
		{
			name: "nil_тело_пишет_только_статус",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				WriteResp(w, http.StatusCreated, nil)
			},
			wantStatus:      http.StatusCreated,
			wantBody:        "",
			wantContentType: false,
		},
		{
			name: "ошибка_кодирования_json_возвращает_500",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				s := struct {
					Ch chan int `json:"ch"`
				}{
					Ch: make(chan int),
				}
				WriteResp(w, http.StatusOK, s)
			},
			wantStatus:      http.StatusInternalServerError,
			wantBody:        `{"title":"SERVER_ERROR","status":500,"detail":"что-то пошло не так"}`,
			wantContentType: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			tt.handler.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())

			resp := rec.Result()
			h := resp.Header.Values("Content-Type")
			if tt.wantContentType {
				require.NotEmpty(t, h)
				require.Contains(t, h[0], "application/json")
				require.Contains(t, h[0], "charset=UTF-8")
				return
			}
			require.Empty(t, h)
		})
	}
}

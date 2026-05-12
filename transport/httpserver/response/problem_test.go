package response

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		title      string
		detail     string
		wantTitle  string
		wantDetail string
	}{
		{
			name:       "not_found_получает_title_по_умолчанию",
			status:     http.StatusNotFound,
			wantTitle:  "NOT_FOUND",
			wantDetail: "",
		},
		{
			name:       "bad_request_получает_title_по_умолчанию",
			status:     http.StatusBadRequest,
			wantTitle:  "VALIDATION_ERROR",
			wantDetail: "",
		},
		{
			name:       "conflict_получает_title_по_умолчанию",
			status:     http.StatusConflict,
			wantTitle:  "CONFLICT",
			wantDetail: "",
		},
		{
			name:       "unprocessable_entity_получает_title_по_умолчанию",
			status:     http.StatusUnprocessableEntity,
			wantTitle:  "BUSINESS_VALIDATION_ERROR",
			wantDetail: "",
		},
		{
			name:       "unsupported_media_type_получает_title_по_умолчанию",
			status:     http.StatusUnsupportedMediaType,
			wantTitle:  "UNSUPPORTED_MEDIA_TYPE",
			wantDetail: "",
		},
		{
			name:       "request_entity_too_large_получает_title_по_умолчанию",
			status:     http.StatusRequestEntityTooLarge,
			wantTitle:  "REQUEST_ENTITY_TOO_LARGE",
			wantDetail: "",
		},
		{
			name:       "неизвестный_для_пакета_статус_получает_http_status_text",
			status:     http.StatusTeapot,
			wantTitle:  http.StatusText(http.StatusTeapot),
			wantDetail: "",
		},
		{
			name:       "кастомный_title_сохраняется",
			status:     http.StatusBadRequest,
			title:      "CUSTOM_ERROR",
			detail:     "детали ошибки",
			wantTitle:  "CUSTOM_ERROR",
			wantDetail: "детали ошибки",
		},
		{
			name:       "internal_server_error_скрывает_detail",
			status:     http.StatusInternalServerError,
			detail:     "техническая причина",
			wantTitle:  "SERVER_ERROR",
			wantDetail: "что-то пошло не так",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Problem(tt.status, tt.title, tt.detail)

			require.Equal(t, tt.status, got.Status)
			require.Equal(t, tt.wantTitle, got.Title)
			require.Equal(t, tt.wantDetail, got.Detail)
		})
	}
}

package response

import (
	"encoding/json"
	"net/http"
)

// WriteResp записывает data как JSON-ответ с указанным HTTP-статусом.
// Если data не удалось закодировать, записывает общую внутреннюю ошибку сервера.
func WriteResp(w http.ResponseWriter, status int, data any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		body, _ := json.Marshal(Problem(http.StatusInternalServerError, "", ""))
		_, _ = w.Write(body)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}

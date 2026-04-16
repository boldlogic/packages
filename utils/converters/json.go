package converters

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrWrongJSON = errors.New("некорректный формат JSON")
)

// DecodeJSON декодирует JSON-данные в структуру типа T.
// Возвращает ошибку, если данные имеют некорректный формат.
func DecodeJSON[T any](filebody []byte) (T, error) {
	var v T
	err := json.Unmarshal(filebody, &v)
	if err != nil {
		return v, fmt.Errorf("%w: %w", ErrWrongJSON, err)
	}
	return v, nil
}

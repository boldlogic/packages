package converters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// DecodeJSONStrict декодирует JSON и запрещает неизвестные поля.
func DecodeJSONStrict[T any](filebody []byte) (T, error) {
	var v T

	decoder := json.NewDecoder(bytes.NewReader(filebody))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&v); err != nil {
		return v, fmt.Errorf("%w: %w", ErrWrongJSON, err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("обнаружено несколько JSON-значений")
		}
		return v, fmt.Errorf("%w: %w", ErrWrongJSON, err)
	}

	return v, nil
}

package httputils

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/boldlogic/packages/utils/converters"
)

// MaxRequestBodySize — максимальный размер тела запроса в байтах, который читает
// [readLimitedRequestBody] и [DecodeRequest] (лишнее отсекается, возвращается [ErrRequestEntityTooLarge]).
const MaxRequestBodySize = 64 * 1024

var (
	// ErrReadingBody возвращается, если прочитать тело HTTP-запроса не удалось.
	ErrReadingBody = errors.New("не удалось прочитать тело запроса")
	// ErrUnsupportedMediaType возвращается, если заголовок Content-Type не задаёт JSON.
	ErrUnsupportedMediaType = errors.New("Content-Type должен быть application/json")
	// ErrRequestEntityTooLarge возвращается, если тело длиннее [MaxRequestBodySize].
	ErrRequestEntityTooLarge = errors.New("тело запроса превышает ограничение")
)

// checkContentType проверяет, что Content-Type задаёт подтип application/json
// (допускаются параметры вроде charset после `;`).
func checkContentType(r *http.Request) error {
	ct := strings.TrimSpace(r.Header.Get("Content-Type"))
	if ct == "" {
		return ErrUnsupportedMediaType
	}
	if !strings.HasPrefix(strings.ToLower(ct), "application/json") {
		return ErrUnsupportedMediaType
	}
	return nil
}

// readLimitedRequestBody читает тело запроса с верхней границей [MaxRequestBodySize].
func readLimitedRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, MaxRequestBodySize+1))
	if err != nil {
		return nil, ErrReadingBody
	}
	if len(body) > MaxRequestBodySize {
		return nil, ErrRequestEntityTooLarge
	}
	return body, nil
}

// DecodeRequest читает тело запроса (с лимитом размера), проверяет Content-Type на
// application/json, декодирует JSON в T без лишних полей и прогоняет go-playground/validator
// по структуре (см. [ValidateStruct]).
func DecodeRequest[T any](r *http.Request) (T, error) {
	var v T
	if err := checkContentType(r); err != nil {
		return v, err
	}
	body, err := readLimitedRequestBody(r)
	if err != nil {
		return v, err
	}

	v, err = converters.DecodeJSONStrict[T](body)
	if err != nil {
		return v, err
	}
	if err := defaultValidator.Struct(v); err != nil {
		return v, err
	}
	return v, nil
}

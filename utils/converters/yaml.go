package converters

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	ErrWrongYAML = errors.New("некорректный формат YAML")
)

// DecodeYAML декодирует YAML-данные в структуру типа T.
// Возвращает ошибку, если данные имеют некорректный формат.
func DecodeYAML[T any](filebody []byte) (T, error) {
	var v T
	err := yaml.Unmarshal(filebody, &v)
	if err != nil {
		return v, fmt.Errorf("%w: %w", ErrWrongYAML, err)
	}
	return v, nil

}

// DecodeYAMLStrict декодирует YAML и запрещает неизвестные поля.
func DecodeYAMLStrict[T any](filebody []byte) (T, error) {
	var v T

	decoder := yaml.NewDecoder(bytes.NewReader(filebody))
	decoder.KnownFields(true)

	if err := decoder.Decode(&v); err != nil {
		if err == io.EOF {
			return v, nil
		}
		return v, fmt.Errorf("%w: %w", ErrWrongYAML, err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("обнаружено несколько YAML-документов")
		}
		return v, fmt.Errorf("%w: %w", ErrWrongYAML, err)
	}

	return v, nil
}

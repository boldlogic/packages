package commonconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boldlogic/packages/utils/converters"
)

var (
	ErrFileStat     = errors.New("не удалось получить информацию о файле")
	ErrFileRead     = errors.New("не удалось прочитать файл")
	ErrWrongFileExt = errors.New("некорректное расширение файла конфигурации")
)

// readFile читает содержимое файла по указанному пути.
// Возвращает ошибку, если файл не существует или недоступен для чтения.
func readFile(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrFileStat, path, err)
	}

	bs, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrFileRead, path, err)
	}
	return bs, nil

}

// DecodeConfig декодирует файл конфигурации в структуру типа T (дженерик).
// Поддерживает форматы: .yaml, .yml, .json.
// Выбор формата определяется по расширению файла.
func DecodeConfig[T any](path string) (T, error) {
	var v T

	fileBody, err := readFile(path)
	if err != nil {
		return v, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch {
	case ext == ".yaml" || ext == ".yml":
		v, err = converters.DecodeYAML[T](fileBody)

	case ext == ".json":
		v, err = converters.DecodeJSON[T](fileBody)

	default:
		return v, fmt.Errorf("%w %s", ErrWrongFileExt, path)
	}

	if err != nil {
		return v, err
	}
	return v, nil
}

// DecodeConfigStrict декодирует файл конфигурации как DecodeConfig,
// но дополнительно запрещает неизвестные поля в YAML/JSON.
func DecodeConfigStrict[T any](path string) (T, error) {
	var v T

	fileBody, err := readFile(path)
	if err != nil {
		return v, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch {
	case ext == ".yaml" || ext == ".yml":
		v, err = converters.DecodeYAMLStrict[T](fileBody)

	case ext == ".json":
		v, err = converters.DecodeJSONStrict[T](fileBody)

	default:
		return v, fmt.Errorf("%w %s", ErrWrongFileExt, path)
	}

	if err != nil {
		return v, err
	}
	return v, nil
}

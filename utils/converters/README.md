# converters

Пакет `converters` предоставляет универсальные функции для декодирования данных в форматах JSON и YAML в типизированные структуры Go.

## Публичный API

```go
func DecodeJSON[T any](filebody []byte) (T, error)
func DecodeJSONStrict[T any](body []byte) (T, error)
func DecodeYAML[T any](filebody []byte) (T, error)
func DecodeYAMLStrict[T any](filebody []byte) (T, error)
```

## Функции

### `DecodeJSON`

Декодирует JSON-данные в структуру типа `T`.

```go
type Config struct {
    Name string `json:"name"`
    Port int    `json:"port"`
}

data := []byte(`{"name": "app", "port": 8080}`)
cfg, err := DecodeJSON[Config](data)
```

### `DecodeJSONStrict`

Декодирует JSON-данные в структуру типа `T`, запрещает неизвестные поля и возвращает ошибку при нескольких JSON-значениях в одном `[]byte`.

### `DecodeYAML`

Декодирует YAML-данные в структуру типа `T`.

```go
type Config struct {
    Name string `yaml:"name"`
    Port int    `yaml:"port"`
}

data := []byte(`
name: app
port: 8080
`)
cfg, err := DecodeYAML[Config](data)
```

### `DecodeYAMLStrict`

Декодирует YAML-данные в структуру типа `T`, запрещает неизвестные поля и возвращает ошибку при нескольких YAML-документах.

## Ошибки

- `ErrWrongJSON` — некорректный формат JSON
- `ErrWrongYAML` — некорректный формат YAML

## Установка

```bash
go get github.com/boldlogic/packages/utils/converters
```

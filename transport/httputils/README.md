# httputils

`httputils` — вспомогательные функции для HTTP-обработчиков: разбор пагинации из query-параметров, чтение JSON-тела запроса с ограничением размера и валидация структур.

## Публичный API

```go
const DefaultLimit = 50
const MaxLimit = 200
const MaxRequestBodySize = 64 * 1024

var ErrInvalidLimit error
var ErrInvalidOffset error
var ErrReadingBody error
var ErrUnsupportedMediaType error
var ErrRequestEntityTooLarge error

func ParseListPagination(r *http.Request) (limit, offset int, err error)
func DecodeRequest[T any](r *http.Request) (T, error)
func ValidateStruct(req any) error
```

## Что делает пакет

- читает `limit` и `offset` из query-параметров HTTP-запроса
- применяет `DefaultLimit`, если `limit` не передан
- ограничивает `limit` значением `MaxLimit`
- проверяет `Content-Type` на `application/json`
- читает тело запроса не больше `MaxRequestBodySize`
- декодирует JSON без лишних полей через strict-декодер
- валидирует структуру через `go-playground/validator`

## Пример

```go
limit, offset, err := httputils.ParseListPagination(r)
if err != nil {
	w.WriteHeader(http.StatusBadRequest)
	return
}

_ = limit
_ = offset
```

```go
type Request struct {
	Name string `json:"name" validate:"required"`
}

req, err := httputils.DecodeRequest[Request](r)
if err != nil {
	w.WriteHeader(http.StatusBadRequest)
	return
}

_ = req
```

## Что важно знать

- `limit` должен быть целым числом больше нуля
- `offset` должен быть целым числом больше либо равен нулю
- если `limit` больше `MaxLimit`, функция возвращает `MaxLimit` без ошибки
- `DecodeRequest` возвращает ошибку, если `Content-Type` не начинается с `application/json`
- `DecodeRequest` запрещает неизвестные поля JSON
- имена полей в ошибках validator берутся из тега `json`

## Установка

```bash
go get github.com/boldlogic/packages/transport/httputils
```

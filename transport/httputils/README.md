# httputils

`httputils` — вспомогательные функции для HTTP-обработчиков: разбор пагинации из query-параметров, чтение JSON-тела запроса с ограничением размера и валидация структур.

## Что делает пакет

- читает `limit` и `offset` из query-параметров HTTP-запроса
- применяет `DefaultLimit`, если `limit` не передан
- ограничивает `limit` значением `MaxLimit`
- проверяет `Content-Type` на `application/json`
- читает тело запроса не больше `MaxRequestBodySize`
- декодирует JSON без лишних полей через strict-декодер
- валидирует структуру через `go-playground/validator`
- регистрирует validator-тег `decimal` для `decimal.Decimal` и строк, разбираемых через `decimal.NewFromString`

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

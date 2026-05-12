# response

`response` - запись JSON-ответов HTTP и описания ошибок в формате, близком к RFC 7807.

## Что делает пакет

- `WriteResp` записывает JSON-ответ с переданным HTTP-статусом
- если тело ответа равно `nil`, `WriteResp` записывает только статус
- если тело не удалось закодировать в JSON, `WriteResp` возвращает `500` с общей ошибкой сервера
- `Problem` формирует структуру ошибки с `title`, `status` и `detail`
- для внутренних ошибок сервера `Problem` скрывает технические детали

## Что важно знать

- `Content-Type: application/json; charset=UTF-8` выставляется только для ответов с телом
- пустой `title` в `Problem` заменяется значением по умолчанию для известного HTTP-статуса
- неизвестные пакету статусы получают `title` из `http.StatusText`

## Установка

```bash
go get github.com/boldlogic/packages/transport/httpserver/response
```

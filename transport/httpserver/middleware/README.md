# middleware

`middleware` — HTTP middleware для debug-логирования запросов и записи HTTP-метрик.

Пакет связывает `zap.Logger` и `httpmetrics.HTTPMetrics`, чтобы оборачивать `http.Handler` без дублирования кода в обработчиках.

## Что делает пакет

- `Wrap` пишет debug-лог после завершения HTTP-запроса
- `Wrap` добавляет в лог метод, путь, query string, статус, заголовки и длительность
- `Metrics` записывает метрики после завершения HTTP-запроса
- `Metrics` использует шаблон маршрута из `chi.RouteContext`, если он доступен
- если шаблон маршрута недоступен, `Metrics` использует `r.URL.Path`
- `Recover` перехватывает panic и пишет ошибку в лог
- если ответ ещё не был записан, `Recover` отправляет статус `500`

## Что важно знать

- `Wrap` и `Metrics` фиксируют статус ответа
- `Wrap` логирует заголовки запроса как строку
- `Metrics` записывает labels `method`, `route` и `status` через пакет `httpmetrics`
- если logger не передан, `NewMiddleware` использует `zap.NewNop()`

## Установка

```bash
go get github.com/boldlogic/packages/transport/httpserver/middleware
```

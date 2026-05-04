# middleware

`middleware` — HTTP middleware для debug-логирования запросов и записи HTTP-метрик.

Пакет связывает `zap.Logger` и `httpmetrics.HTTPMetrics`, чтобы оборачивать `http.Handler` без дублирования кода в обработчиках.

## Публичный API

```go
type Middleware struct {}

func NewMiddleware(metrics *httpmetrics.HTTPMetrics, logger *zap.Logger) *Middleware
func (m Middleware) Wrap(next http.Handler) http.Handler
func (m Middleware) Metrics(next http.Handler) http.Handler
```

## Что делает пакет

- `Wrap` пишет debug-лог после завершения HTTP-запроса
- `Wrap` добавляет в лог метод, путь, query string, статус, заголовки и длительность
- `Metrics` записывает метрики после завершения HTTP-запроса
- `Metrics` использует шаблон маршрута из `chi.RouteContext`, если он доступен
- если шаблон маршрута недоступен, `Metrics` использует `r.URL.Path`
- статус ответа считывается через внутренний `responseRecorder`

## Пример

```go
reg := metrics.New()
httpMetrics := httpmetrics.NewMetrics(reg)
mw := middleware.NewMiddleware(httpMetrics, logger)

handler := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}))

handler = mw.Metrics(handler)
```

## Что важно знать

- `responseRecorder` по умолчанию считает статусом `http.StatusOK`
- `responseRecorder` фиксирует статус, переданный через `WriteHeader`
- `Wrap` логирует заголовки запроса как строку
- `Metrics` записывает labels `method`, `route` и `status` через пакет `httpmetrics`
- `NewMiddleware` сохраняет переданные зависимости без дополнительной инициализации

## Установка

```bash
go get github.com/boldlogic/packages/transport/httpserver/middleware
```

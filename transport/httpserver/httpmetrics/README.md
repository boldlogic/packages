# httpmetrics

`httpmetrics` — набор Prometheus-метрик для HTTP-запросов.

Пакет регистрирует счётчик запросов и гистограмму длительности в переданном registry, а затем позволяет записывать наблюдения через `RecordRequest`.

## Публичный API

```go
type HTTPMetrics struct {}

func NewMetrics(reg metrics.Registry) *HTTPMetrics
func (m *HTTPMetrics) RecordRequest(method, route, status string, duration time.Duration)
```

## Метрики

```text
http_requests_total{method,route,status}
http_request_duration_seconds{method,route}
```

## Что делает пакет

- создаёт `CounterVec` для общего числа HTTP-запросов
- создаёт `HistogramVec` для длительности HTTP-запросов в секундах
- регистрирует метрики в переданном registry
- переиспользует уже зарегистрированные collector, если registry вернул `prometheus.AlreadyRegisteredError`
- записывает метод, маршрут, статус и длительность запроса

## Пример

```go
reg := metrics.New()
httpMetrics := httpmetrics.NewMetrics(reg)

httpMetrics.RecordRequest("GET", "/health", "200", 10*time.Millisecond)
```

## Установка

```bash
go get github.com/boldlogic/packages/transport/httpserver/httpmetrics
```


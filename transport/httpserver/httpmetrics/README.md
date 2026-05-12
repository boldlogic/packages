# httpmetrics

`httpmetrics` — набор Prometheus-метрик для HTTP-запросов.

Пакет регистрирует счётчик запросов и гистограмму длительности в переданном registry, а затем позволяет записывать наблюдения через `RecordRequest`.

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

## Установка

```bash
go get github.com/boldlogic/packages/transport/httpserver/httpmetrics
```


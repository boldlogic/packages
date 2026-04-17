# metrics

`metrics` — небольшой пакет для создания базового Prometheus registry с уже подключёнными стандартными коллекторами.

Пакет полезен, когда в сервисе нужен единый способ поднять registry без повторения boilerplate-кода.

## Публичный API

```go
type Registry interface {
    prometheus.Registerer
    prometheus.Gatherer
}

func New() *prometheus.Registry
```

## Что делает пакет

- создаёт новый `prometheus.Registry`
- регистрирует `GoCollector`
- регистрирует `ProcessCollector`

## Пример

```go
reg := metrics.New()

counter := prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests_total",
	Help: "Количество запросов",
})

if err := reg.Register(counter); err != nil {
	log.Printf("не удалось зарегистрировать метрику: %v", err)
}
```

## Что важно знать

- пакет не экспортирует собственную реализацию registry, а возвращает `*prometheus.Registry`
- стандартные коллекторы подключаются автоматически при вызове `New`

## Установка

```bash
go get github.com/boldlogic/packages/metrics
```

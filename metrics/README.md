# metrics

`metrics` — небольшой пакет для создания базового Prometheus registry с уже подключёнными стандартными коллекторами.

## Что делает пакет

- создаёт новый `prometheus.Registry`
- регистрирует `GoCollector`
- регистрирует `ProcessCollector`

## Что важно знать

- пакет не экспортирует собственную реализацию registry, а возвращает `*prometheus.Registry`
- стандартные коллекторы подключаются автоматически при вызове `New`

## Установка

```bash
go get github.com/boldlogic/packages/metrics
```

# periodic

`periodic` — пакет для запуска периодических задач и координации нескольких worker через общий `context.Context`.

Пакет полезен для фоновых джобов в сервисах: синхронизация, обновление кэша, polling, housekeeping-задачи и другие процессы, которые нужно выполнять по интервалу.

## Публичный API

```go
type Worker interface {
    Name() string
    Run(ctx context.Context)
}

type Job func(ctx context.Context) error

type PeriodicWorker struct {}

func NewPeriodicWorker(name, errMsg string, interval time.Duration, job Job, logger *zap.Logger) *PeriodicWorker
func (w *PeriodicWorker) Name() string
func (w *PeriodicWorker) Run(ctx context.Context)

type Runner struct {}

func NewRunner(workers ...Worker) *Runner
func (r *Runner) Run(ctx context.Context)
```

## Что делает пакет

- запускает `Job` по заданному интервалу
- завершает worker при отмене контекста
- логирует ошибку выполнения через переданный `zap.Logger`
- умеет запускать несколько worker параллельно через `Runner`

## Пример

```go
worker := periodic.NewPeriodicWorker(
	"refresh-cache",
	"не удалось обновить кэш",
	time.Minute,
	func(ctx context.Context) error {
		return refreshCache(ctx)
	},
	logger,
)

runner := periodic.NewRunner(worker)
runner.Run(ctx)
```

## Что важно знать

- если интервал меньше либо равен нулю, используется `60 * time.Second`
- после ошибки `PeriodicWorker` ждёт `5 * time.Second` перед следующей попыткой
- `Runner.Run` блокируется до отмены контекста и завершения всех worker
- пакет не управляет восстановлением после panic внутри `Job`

## Установка

```bash
go get github.com/boldlogic/packages/periodic
```

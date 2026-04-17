package periodic

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Worker описывает задачу, которую можно запускать под управлением Runner.
type Worker interface {
	Name() string
	Run(ctx context.Context)
}

// Job описывает функцию периодической работы, которая выполняется по расписанию.
type Job func(ctx context.Context) error

// PeriodicWorker выполняет Job с заданным интервалом и логирует ошибки выполнения.
type PeriodicWorker struct {
	name       string
	interval   time.Duration
	retryDelay time.Duration
	job        Job
	logger     *zap.Logger
	errMsg     string
}

// NewPeriodicWorker создаёт новый PeriodicWorker.
//
// Если интервал не задан или меньше либо равен нулю, используется значение
// по умолчанию `60 * time.Second`.
func NewPeriodicWorker(name, errMsg string, interval time.Duration, job Job, logger *zap.Logger) *PeriodicWorker {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &PeriodicWorker{
		name:       name,
		interval:   interval,
		retryDelay: 5 * time.Second,
		job:        job,
		logger:     logger,
		errMsg:     errMsg,
	}
}

// Name возвращает имя worker.
func (w *PeriodicWorker) Name() string { return w.name }

// Run запускает периодическое выполнение Job до отмены контекста.
func (w *PeriodicWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.job(ctx)
			if err != nil {
				w.logger.Error(w.errMsg, zap.Error(err))
				select {
				case <-time.After(w.retryDelay):
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// Runner управляет запуском набора Worker и ожидает их завершения.
type Runner struct {
	workers []Worker
	wg      sync.WaitGroup
}

// NewRunner создаёт новый Runner для переданных Worker.
func NewRunner(workers ...Worker) *Runner {
	return &Runner{workers: workers}
}

// Run запускает все Worker и блокируется до отмены контекста и завершения всех горутин.
func (r *Runner) Run(ctx context.Context) {
	for _, w := range r.workers {
		w := w
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			w.Run(ctx)
		}()
	}
	<-ctx.Done()
	r.wg.Wait()
}

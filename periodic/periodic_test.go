package periodic

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"
)

type stubWorker struct {
	name string
	run  func(ctx context.Context)
}

func (w stubWorker) Name() string { return w.name }

func (w stubWorker) Run(ctx context.Context) {
	if w.run != nil {
		w.run(ctx)
	}
}

func TestNewPeriodicWorker(t *testing.T) {
	t.Run("интервал_по_умолчанию", func(t *testing.T) {
		w := NewPeriodicWorker("job", "ошибка job", 0, func(_ context.Context) error { return nil }, zap.NewNop())
		if w.interval != 60*time.Second {
			t.Fatalf("получили %v, ожидали %v", w.interval, 60*time.Second)
		}
	})

	t.Run("заданный_интервал", func(t *testing.T) {
		w := NewPeriodicWorker("job", "ошибка job", time.Second, func(_ context.Context) error { return nil }, zap.NewNop())
		if w.interval != time.Second {
			t.Fatalf("получили %v, ожидали %v", w.interval, time.Second)
		}
	})
}

func TestPeriodicWorker_Name(t *testing.T) {
	t.Run("возвращает_имя", func(t *testing.T) {
		w := NewPeriodicWorker("job", "ошибка job", time.Second, func(_ context.Context) error { return nil }, zap.NewNop())
		if w.Name() != "job" {
			t.Fatalf("получили %q, ожидали job", w.Name())
		}
	})
}

func TestPeriodicWorker_Run(t *testing.T) {
	t.Run("выполнение_до_отмены_контекста", func(t *testing.T) {
		var called atomic.Int32
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		w := NewPeriodicWorker("job", "ошибка job", 10*time.Millisecond, func(_ context.Context) error {
			if called.Add(1) == 1 {
				cancel()
			}
			return nil
		}, zap.NewNop())

		done := make(chan struct{})
		go func() {
			w.Run(ctx)
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("ожидали завершение worker")
		}

		if called.Load() == 0 {
			t.Fatal("ожидали хотя бы один вызов job")
		}
	})

	t.Run("ошибка_job", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var called atomic.Int32
		w := NewPeriodicWorker("job", "ошибка job", 10*time.Millisecond, func(_ context.Context) error {
			if called.Add(1) == 1 {
				cancel()
			}
			return errors.New("ошибка выполнения")
		}, zap.NewNop())
		w.retryDelay = 1 * time.Millisecond

		done := make(chan struct{})
		go func() {
			w.Run(ctx)
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("ожидали завершение worker")
		}

		if called.Load() == 0 {
			t.Fatal("ожидали хотя бы один вызов job")
		}
	})
}

func TestNewRunner(t *testing.T) {
	t.Run("создание_runner", func(t *testing.T) {
		r := NewRunner()
		if r == nil {
			t.Fatal("ожидали runner")
		}
	})
}

func TestRunner_Run(t *testing.T) {
	t.Run("запуск_всех_worker", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var called atomic.Int32
		r := NewRunner(
			stubWorker{name: "one", run: func(ctx context.Context) { called.Add(1); <-ctx.Done() }},
			stubWorker{name: "two", run: func(ctx context.Context) { called.Add(1); <-ctx.Done() }},
		)

		done := make(chan struct{})
		go func() {
			r.Run(ctx)
			close(done)
		}()

		deadline := time.Now().Add(time.Second)
		for called.Load() < 2 && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}

		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("ожидали завершение runner")
		}

		if called.Load() != 2 {
			t.Fatalf("получили %d запусков, ожидали 2", called.Load())
		}
	})
}

package shutdown

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestIsExceeded(t *testing.T) {
	t.Run("ошибка_nil", func(t *testing.T) {
		if IsExceeded(nil) {
			t.Fatal("ожидали false")
		}
	})

	t.Run("context_Canceled", func(t *testing.T) {
		if !IsExceeded(context.Canceled) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("обёртка_context_Canceled", func(t *testing.T) {
		err := fmt.Errorf("ошибка БД: %w", context.Canceled)
		if !IsExceeded(err) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("context_DeadlineExceeded", func(t *testing.T) {
		if !IsExceeded(context.DeadlineExceeded) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("обёртка_context_DeadlineExceeded", func(t *testing.T) {
		err := fmt.Errorf("превышено время ожидания: %w", context.DeadlineExceeded)
		if !IsExceeded(err) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("другая_ошибка", func(t *testing.T) {
		if IsExceeded(errors.New("произвольная ошибка")) {
			t.Fatal("ожидали false")
		}
	})
}

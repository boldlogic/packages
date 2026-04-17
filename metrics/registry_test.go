package metrics

import "testing"

func TestNew(t *testing.T) {
	t.Run("создание_registry", func(t *testing.T) {
		reg := New()
		if reg == nil {
			t.Fatal("ожидали registry")
		}
	})

	t.Run("go_и_process_коллекторы_зарегистрированы", func(t *testing.T) {
		reg := New()

		metrics, err := reg.Gather()
		if err != nil {
			t.Fatalf("сбор метрик: %v", err)
		}
		if len(metrics) == 0 {
			t.Fatal("ожидали зарегистрированные метрики")
		}
	})
}

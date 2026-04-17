package logger

import "testing"

func TestConfig_Validate(t *testing.T) {
	t.Run("level_info", func(t *testing.T) {
		errs := (Config{Level: "info"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("level_debug", func(t *testing.T) {
		errs := (Config{Level: "debug"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("level_error", func(t *testing.T) {
		errs := (Config{Level: "error"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("level_warn", func(t *testing.T) {
		errs := (Config{Level: "warn"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})
	t.Run("level_panic", func(t *testing.T) {
		errs := (Config{Level: "panic"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("level_не_парсится", func(t *testing.T) {
		errs := (Config{Level: "не-уровень"}).Validate()
		if len(errs) != 1 {
			t.Fatalf("ожидали 1 ошибку, получили %v", errs)
		}
	})
}

package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_writeSyncers(t *testing.T) {
	t.Run("output_file_пуст", func(t *testing.T) {
		ss := writeSyncers(Config{OutputFile: ""})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})

	t.Run("output_file_открывается", func(t *testing.T) {
		ss := writeSyncers(Config{OutputFile: os.DevNull})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})

	t.Run("output_file_недоступен", func(t *testing.T) {
		badDir := filepath.Join(t.TempDir(), "нет_каталога", "x.log")
		ss := writeSyncers(Config{OutputFile: badDir})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})

	t.Run("output_file_переиспользуется", func(t *testing.T) {
		first := writeSyncers(Config{OutputFile: os.DevNull})
		second := writeSyncers(Config{OutputFile: os.DevNull})

		if len(first) != 1 || len(second) != 1 {
			t.Fatalf("ожидали по одному syncer, получили %d и %d", len(first), len(second))
		}
		if first[0] != second[0] {
			t.Fatal("ожидали переиспользование syncer для одного файла")
		}
	})
}

func TestNew_noPanic(t *testing.T) {
	t.Run("new_output_file_пуст", func(t *testing.T) {
		log := New(Config{Level: "info", Format: "console", OutputFile: ""})
		if log == nil {
			t.Fatal("ожидали logger")
		}
		_ = log.Sync()
	})

	t.Run("new_output_file_devnull", func(t *testing.T) {
		log := New(Config{Level: "info", Format: "json", OutputFile: os.DevNull})
		if log == nil {
			t.Fatal("ожидали logger")
		}
		log.Info("test")
		_ = log.Sync()
	})
}

package commonconfig

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	t.Run("возвращает_дефолтный_путь_без_флага", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		os.Args = []string{"test"}

		defaultPath := filepath.Join("config", "default.yaml")
		got := GetConfigPath(defaultPath)

		if got != defaultPath {
			t.Fatalf("ожидали %q, получили %q", defaultPath, got)
		}
	})

	t.Run("возвращает_путь_из_флага", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		customPath := filepath.Join("custom", "config.yaml")
		os.Args = []string{"test", "-config", customPath}

		got := GetConfigPath("default.yaml")

		if got != customPath {
			t.Fatalf("ожидали %q, получили %q", customPath, got)
		}
	})

}

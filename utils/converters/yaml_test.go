package converters

import (
	"errors"
	"testing"
)

type sampleYAML struct {
	Name string `yaml:"name"`
}

func TestDecodeYAML(t *testing.T) {
	t.Run("валидный_yaml", func(t *testing.T) {
		got, err := DecodeYAML[sampleYAML]([]byte("name: one\n"))
		if err != nil {
			t.Fatalf("DecodeYAML: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("обрыв_yaml_ErrWrongYAML", func(t *testing.T) {
		_, err := DecodeYAML[sampleYAML]([]byte("name: [\n"))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongYAML) {
			t.Fatalf("ошибка %v, ожидали ErrWrongYAML", err)
		}
	})

	t.Run("пустой_yaml", func(t *testing.T) {
		got, err := DecodeYAML[sampleYAML]([]byte(""))
		if err != nil {
			t.Fatalf("DecodeYAML: %v", err)
		}
		if got.Name != "" {
			t.Fatalf("Name=%q, ожидали пусто", got.Name)
		}
	})
}

func TestDecodeYAMLStrict(t *testing.T) {
	t.Run("strict_валидный_yaml", func(t *testing.T) {
		got, err := DecodeYAMLStrict[sampleYAML]([]byte("name: one\n"))
		if err != nil {
			t.Fatalf("DecodeYAMLStrict: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("strict_неизвестное_поле_ErrWrongYAML", func(t *testing.T) {
		_, err := DecodeYAMLStrict[sampleYAML]([]byte("name: one\nextra: value\n"))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongYAML) {
			t.Fatalf("ошибка %v, ожидали ErrWrongYAML", err)
		}
	})
}

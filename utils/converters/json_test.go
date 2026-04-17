package converters

import (
	"errors"
	"testing"
)

type sampleJSON struct {
	Name string `json:"name"`
}

func TestDecodeJSON(t *testing.T) {
	t.Run("валидный_json", func(t *testing.T) {
		got, err := DecodeJSON[sampleJSON]([]byte(`{"name":"one"}`))
		if err != nil {
			t.Fatalf("DecodeJSON: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("обрыв_json_ErrWrongJSON", func(t *testing.T) {
		_, err := DecodeJSON[sampleJSON]([]byte(`{"name":`))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongJSON) {
			t.Fatalf("ошибка %v, ожидали ErrWrongJSON", err)
		}
	})

	t.Run("пустое_тело_ErrWrongJSON", func(t *testing.T) {
		_, err := DecodeJSON[sampleJSON]([]byte(""))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongJSON) {
			t.Fatalf("ошибка %v, ожидали ErrWrongJSON", err)
		}
	})
}

func TestDecodeJSONStrict(t *testing.T) {
	t.Run("strict_валидный_json", func(t *testing.T) {
		got, err := DecodeJSONStrict[sampleJSON]([]byte(`{"name":"one"}`))
		if err != nil {
			t.Fatalf("DecodeJSONStrict: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("strict_неизвестное_поле_ErrWrongJSON", func(t *testing.T) {
		_, err := DecodeJSONStrict[sampleJSON]([]byte(`{"name":"one","extra":"value"}`))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongJSON) {
			t.Fatalf("ошибка %v, ожидали ErrWrongJSON", err)
		}
	})
}

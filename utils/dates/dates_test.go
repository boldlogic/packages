package dates

import (
	"errors"
	"testing"
)

func TestOptionalDatePtr(t *testing.T) {
	t.Run("передана_корректная_строка_даты_с_ISODateFormat", func(t *testing.T) {
		date := "2024-06-19"
		got, err := OptionalDatePtr(date, ISODateFormat)
		if err != nil {
			t.Fatalf("OptionalDatePtr: %v", err)
		}
		if got.Year() != 2024 || got.Month() != 6 || got.Day() != 19 {
			t.Fatalf("ожидали 2024-06-19, получили %v", got.Format(ISODateFormat))
		}
	})
	t.Run("передана_пустая_строка_с_ISODateFormat", func(t *testing.T) {
		date := ""
		got, err := OptionalDatePtr(date, ISODateFormat)
		if err != nil {
			t.Fatalf("OptionalDatePtr: %v", err)
		}
		if got != nil {
			t.Fatalf("ожидали nil, получили %v", got.Format(ISODateFormat))
		}
	})
	t.Run("передана_некорректная_строка_даты_с_ISODateFormat", func(t *testing.T) {
		date := "2024-06"
		got, err := OptionalDatePtr(date, ISODateFormat)
		if !errors.Is(err, ErrWrongDateFormat) {
			t.Fatalf("OptionalDatePtr: %v", err)
		}
		if got != nil {
			t.Fatalf("ожидали nil, получили %v", got.Format(ISODateFormat))
		}
	})
	t.Run("передана_корректная_строка_даты_с_кастомным_форматом", func(t *testing.T) {
		date := "2024-06-19"
		got, err := OptionalDatePtr(date, "2006-01-02")
		if err != nil {
			t.Fatalf("OptionalDatePtr: %v", err)
		}
		if got.Year() != 2024 || got.Month() != 6 || got.Day() != 19 {
			t.Fatalf("ожидали 2024-06-19, получили %v", got.Format(ISODateFormat))
		}
	})
}

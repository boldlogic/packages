package httputils

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_isDecimal(t *testing.T) {

	t.Run("корректный_decimal", func(t *testing.T) {

		s := struct {
			D decimal.Decimal `json:"balance" validate:"decimal"`
		}{
			D: decimal.RequireFromString("11.01"),
		}
		err := defaultValidator.Struct(s)
		if err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("корректный_decimalPtr", func(t *testing.T) {
		d := decimal.RequireFromString("11.01")
		s := struct {
			D *decimal.Decimal `json:"balance" validate:"required,decimal"`
		}{
			D: &d,
		}
		err := defaultValidator.Struct(s)
		if err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("корректный_string", func(t *testing.T) {

		s := struct {
			D string `json:"balance" validate:"decimal"`
		}{
			D: "1101",
		}
		err := defaultValidator.Struct(s)
		if err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("пустой_decimal", func(t *testing.T) {
		s := struct {
			D *decimal.Decimal `json:"balance" validate:"required,decimal"`
		}{}
		err := defaultValidator.Struct(s)
		if err == nil {
			t.Fatalf("ожидали ошибку")
		}
	})
	t.Run("пустой_decimal_ptr", func(t *testing.T) {
		s := struct {
			D *decimal.Decimal `json:"balance" validate:"omitempty,decimal"`
		}{}
		err := defaultValidator.Struct(s)
		if err != nil {
			t.Fatalf("ожидали ошибку")
		}
	})
	t.Run("нет_строки", func(t *testing.T) {

		s := struct {
			D string `json:"balance" validate:"decimal"`
		}{}
		err := defaultValidator.Struct(s)
		if err == nil {
			t.Fatalf("%v", err)
		}
	})
	t.Run("нет_строки_ptr", func(t *testing.T) {

		s := struct {
			D *string `json:"balance" validate:"decimal"`
		}{}
		err := defaultValidator.Struct(s)
		if err == nil {
			t.Fatalf("ожидали ошибку")
		}
	})

}

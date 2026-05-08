package httputils

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func isDecimal(fl validator.FieldLevel) bool {
	switch v := fl.Field().Interface().(type) {
	case decimal.Decimal:
		return true
	// case *decimal.Decimal:
	// 	return v != nil
	case string:
		_, err := decimal.NewFromString(v)

		return err == nil
	default:
		return false
	}
}

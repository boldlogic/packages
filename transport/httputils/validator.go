package httputils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// defaultValidator — общий инстанс для [DecodeRequest] и [ValidateStruct]
// (имена полей в ошибках берутся из тега json).
var defaultValidator *validator.Validate

func init() {
	defaultValidator = validator.New()

	defaultValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct проверяет значение через go-playground/validator с теми же
// настройками тегов, что и при разборе JSON из запроса.
func ValidateStruct(req any) error {
	return defaultValidator.Struct(req)
}

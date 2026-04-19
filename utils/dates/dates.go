package dates

import (
	"errors"
	"time"
)

const (
	ISODateFormat string = "2006-01-02" // ISO 8601 (YYYY-MM-DD)
)

var (
	ErrWrongDateFormat = errors.New("некорректный формат date")
)

// OptionalDatePtr преобразует строку даты в `*time.Time`, если дата передана.
//
// Если входная строка пуста, функция возвращает `nil, nil`.
// Если строка задана и соответствует переданному layout, функция возвращает
// указатель на распарсенное значение времени.
// Если строка не соответствует формату, возвращается `ErrWrongDateFormat`.
func OptionalDatePtr(date string, layout string) (*time.Time, error) {
	var dt *time.Time
	if date != "" {

		parsed, err := time.Parse(layout, date)
		if err != nil {
			return nil, ErrWrongDateFormat
		}
		dt = &parsed
	}
	return dt, nil
}

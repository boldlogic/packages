package dates

import (
	"errors"
	"fmt"
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

	if date == "" {
		return dt, nil
	}

	parsed, err := time.Parse(layout, date)
	if err != nil {
		return nil, ErrWrongDateFormat
	}

	return &parsed, nil
}

// ParseWithDefaultNow парсит строку date с помощью time.Parse и переданного layout.
//
// Если date пустая, возвращается time.Now и nil — удобно для значений по умолчанию
// из query-параметров или форм.
// Если строка не соответствует формату, ошибка обёртывает ErrWrongDateFormat
// (errors.Is(err, ErrWrongDateFormat) остаётся истинным); в цепочке сохраняется
// сообщение от time.Parse для отладки.
func ParseWithDefaultNow(date string, layout string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil
	}

	parsed, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrWrongDateFormat, err)
	}

	return parsed, nil

}

// Today Возвращает сегодняшнюю дату (полночь) в локальной временной зоне.
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TruncateToDateIn Обрезает t до даты (полночь) в указанной локали.
func TruncateToDateIn(t time.Time, loc *time.Location) time.Time {
	tIn := t.In(loc)
	return time.Date(tIn.Year(), tIn.Month(), tIn.Day(), 0, 0, 0, 0, loc)
}

// DateToYYYYMMDD Возвращает дату как int64 YYYYMMDD
func DateToYYYYMMDD(t time.Time) int64 {
	y, m, d := t.Date()
	return int64(y)*10000 + int64(m)*100 + int64(d)
}

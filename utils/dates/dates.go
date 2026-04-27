package dates

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ISODateFormat string = "2006-01-02" // ISO 8601 (YYYY-MM-DD)
	// DatetimeFormat — дата и время «наивные»; в [ParseScheduledAt] разбираются в UTC.
	DatetimeFormat = "2006-01-02 15:04:05"
)

var (
	ErrWrongDateFormat        = errors.New("некорректный формат date")
	ErrWrongScheduledAtFormat = errors.New("некорректный формат scheduledAt. Ожидается RFC3339")
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

// EarliestDate возвращает указатель на наименьшую (самую раннюю) дату среди dates.
// Элементы nil пропускаются. Если не передано ни одной ненилевой даты, возвращается nil.
// Возвращается один из исходных указателей, а не копия [time.Time].
func EarliestDate(dates ...*time.Time) *time.Time {
	var result *time.Time
	for _, d := range dates {
		if d == nil {
			continue
		}
		if result == nil || d.Before(*result) {
			result = d
		}
	}
	return result
}

// ParseScheduledAt разбирает момент времени для планирования (scheduledAt).
//
// Строка из одних пробелов и табуляций считается пустой: возвращается time.Now и nil.
// Иначе по очереди пробуются:
//   - time.RFC3339Nano;
//   - time.RFC3339;
//   - [DatetimeFormat] через time.ParseInLocation(…, time.UTC).
//
// Если ни один вариант не подошёл — (time.Time{}, [ErrWrongScheduledAtFormat]).
func ParseScheduledAt(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation(DatetimeFormat, s, time.UTC); err == nil {
		return t, nil
	}
	return time.Time{}, ErrWrongScheduledAtFormat
}

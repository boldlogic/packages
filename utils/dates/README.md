# dates

`dates` — утилиты для разбора, нормализации и сравнения дат.

Пакет содержит общие форматы дат, ошибки разбора и функции для случаев, где пустое значение должно трактоваться как отсутствие даты или как текущий момент.

## Публичный API

```go
const ISODateFormat = "2006-01-02"
const DatetimeFormat = "2006-01-02 15:04:05"

var ErrWrongDateFormat error
var ErrWrongScheduledAtFormat error

func OptionalDatePtr(date string, layout string) (*time.Time, error)
func ParseWithDefaultNow(date string, layout string) (time.Time, error)
func Today() time.Time
func TruncateToDateIn(t time.Time, loc *time.Location) time.Time
func DateToYYYYMMDD(t time.Time) int64
func EarliestDate(dates ...*time.Time) *time.Time
func ParseScheduledAt(s string) (time.Time, error)
```

## Что делает пакет

- парсит необязательную дату в `*time.Time`
- возвращает `time.Now()`, если строка даты не передана
- возвращает текущую дату с временем `00:00:00` в локальной временной зоне
- обрезает `time.Time` до даты в указанной локации
- преобразует дату в число формата `YYYYMMDD`
- выбирает самую раннюю дату из набора указателей
- разбирает `scheduledAt` в форматах `time.RFC3339Nano`, `time.RFC3339` и `DatetimeFormat`

## Пример

```go
date, err := dates.OptionalDatePtr("2026-05-04", dates.ISODateFormat)
if err != nil {
	return err
}

_ = date
```

```go
scheduledAt, err := dates.ParseScheduledAt("2026-05-04T12:30:00Z")
if err != nil {
	return err
}

_ = scheduledAt
```

## Что важно знать

- `OptionalDatePtr` возвращает `nil, nil`, если входная строка пустая
- `ParseWithDefaultNow` возвращает `time.Now()`, если входная строка пустая
- `ParseScheduledAt` считает строку из пробелов и табуляций пустой
- `DatetimeFormat` в `ParseScheduledAt` разбирается в UTC
- `EarliestDate` возвращает один из исходных указателей, а не копию даты

## Установка

```bash
go get github.com/boldlogic/packages/utils/dates
```

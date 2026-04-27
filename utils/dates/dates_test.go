package dates

import (
	"errors"
	"testing"
	"time"
)

func TestOptionalDatePtr(t *testing.T) {
	t.Parallel()
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

func TestParseWithDefaultNow(t *testing.T) {
	t.Parallel()
	t.Run("передана_корректная_строка_даты_с_ISODateFormat", func(t *testing.T) {
		date := "2024-06-19"
		got, err := ParseWithDefaultNow(date, ISODateFormat)
		if err != nil {
			t.Fatalf("ParseWithDefaultNow: %v", err)
		}
		if got.Year() != 2024 || got.Month() != 6 || got.Day() != 19 {
			t.Fatalf("ParseWithDefaultNow: ожидали 2024-06-19, получили %v", got.Format(ISODateFormat))
		}
	})
	t.Run("передана_пустая_строка_с_ISODateFormat", func(t *testing.T) {
		date := ""
		got, err := ParseWithDefaultNow(date, ISODateFormat)
		if err != nil {
			t.Fatalf("ParseWithDefaultNow: %v", err)
		}
		want := time.Now()
		if (got.Year() != want.Year()) || (got.Month() != want.Month()) || (got.Day() != want.Day()) {
			t.Fatalf("ParseWithDefaultNow: ожидали %v, получили %v", want, got.Format(ISODateFormat))
		}
	})
	t.Run("передана_некорректная_строка_даты_с_ISODateFormat", func(t *testing.T) {
		date := "2024-06"
		got, err := ParseWithDefaultNow(date, ISODateFormat)
		if !errors.Is(err, ErrWrongDateFormat) {
			t.Fatalf("ParseWithDefaultNow: %v", err)
		}
		want := time.Time{}
		if got != want {
			t.Fatalf("ParseWithDefaultNow: ожидали %v, получили %v", want, got.Format(ISODateFormat))
		}
	})
	t.Run("передана_корректная_строка_даты_с_кастомным_форматом", func(t *testing.T) {
		date := "2024-06-19"
		got, err := ParseWithDefaultNow(date, "2006-01-02")
		if err != nil {
			t.Fatalf("ParseWithDefaultNow: %v", err)
		}
		if got.Year() != 2024 || got.Month() != 6 || got.Day() != 19 {
			t.Fatalf("ParseWithDefaultNow: ожидали 2024-06-19, получили %v", got.Format(ISODateFormat))
		}
	})
}

func TestToday(t *testing.T) {
	t.Parallel()
	t.Run("совпадает_с_календарной_датой_локали_и_полночь", func(t *testing.T) {
		ref := time.Now()
		got := Today()
		if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
			t.Fatalf("Today: ожидали полночь, получили %v", got)
		}
		if got.Location() != ref.Location() {
			t.Fatalf("Today: ожидали локаль %v, получили %v", ref.Location(), got.Location())
		}
		if got.Year() != ref.Year() || got.Month() != ref.Month() || got.Day() != ref.Day() {
			t.Fatalf("Today: ожидали дату как у time.Now (%v), получили %v", ref.Format(ISODateFormat), got.Format(ISODateFormat))
		}
	})
}

func TestTruncateToDateIn(t *testing.T) {
	t.Parallel()
	t.Run("обрезка_времени_до_полуночи_в_UTC", func(t *testing.T) {
		tm := time.Date(2024, 6, 19, 14, 30, 45, 123456789, time.UTC)
		got := TruncateToDateIn(tm, time.UTC)
		want := time.Date(2024, 6, 19, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Fatalf("TruncateToDateIn: ожидали %v, получили %v", want, got)
		}
	})
	t.Run("дата_в_другой_локали_сохраняет_календарный_день_в_этой_локали", func(t *testing.T) {
		plus2 := time.FixedZone("plus2", 2*3600)
		// 2024-06-19 01:00 UTC = 2024-06-19 03:00 в +2
		tm := time.Date(2024, 6, 19, 1, 0, 0, 0, time.UTC)
		got := TruncateToDateIn(tm, plus2)
		want := time.Date(2024, 6, 19, 0, 0, 0, 0, plus2)
		if !got.Equal(want) {
			t.Fatalf("TruncateToDateIn: ожидали %v, получили %v", want, got)
		}
	})
	t.Run("смена_календарного_дня_при_переводе_в_локаль", func(t *testing.T) {
		plus2 := time.FixedZone("plus2", 2*3600)
		// 2024-06-19 22:00 UTC = 2024-06-20 00:00 в +2 — полночь по календарю +2 это 20 июня
		tm := time.Date(2024, 6, 19, 22, 0, 0, 0, time.UTC)
		got := TruncateToDateIn(tm, plus2)
		want := time.Date(2024, 6, 20, 0, 0, 0, 0, plus2)
		if !got.Equal(want) {
			t.Fatalf("TruncateToDateIn: ожидали %v, получили %v", want, got)
		}
	})
}

func TestDateToYYYYMMDD(t *testing.T) {
	t.Parallel()
	t.Run("обычная_дата", func(t *testing.T) {
		tm := time.Date(2024, 6, 19, 12, 0, 0, 0, time.UTC)
		got := DateToYYYYMMDD(tm)
		const want int64 = 20240619
		if got != want {
			t.Fatalf("DateToYYYYMMDD: ожидали %d, получили %d", want, got)
		}
	})
	t.Run("однозначные_месяц_и_день", func(t *testing.T) {
		tm := time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)
		got := DateToYYYYMMDD(tm)
		const want int64 = 20240305
		if got != want {
			t.Fatalf("DateToYYYYMMDD: ожидали %d, получили %d", want, got)
		}
	})
	t.Run("игнорирует_время_суток", func(t *testing.T) {
		tm := time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC)
		got := DateToYYYYMMDD(tm)
		const want int64 = 20241231
		if got != want {
			t.Fatalf("DateToYYYYMMDD: ожидали %d, получили %d", want, got)
		}
	})
}

func TestEarliestDate(t *testing.T) {
	t.Parallel()
	t.Run("только_nil", func(t *testing.T) {
		if got := EarliestDate(nil, nil); got != nil {
			t.Fatalf("EarliestDate: ожидали nil, получили %v", got)
		}
		if got := EarliestDate(); got != nil {
			t.Fatalf("EarliestDate: ожидали nil, получили %v", got)
		}
	})
	t.Run("один_ненилевой", func(t *testing.T) {
		a := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
		pa := &a
		if got := EarliestDate(nil, pa, nil); got != &a {
			t.Fatalf("EarliestDate: ожидали тот же указатель, got=%v", got)
		}
	})
	t.Run("раньше_другого", func(t *testing.T) {
		early := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		late := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		pe, pl := &early, &late
		if got := EarliestDate(pl, pe); got != pe {
			t.Fatalf("EarliestDate: ожидали %v, получили %v", &early, got)
		}
		if got := EarliestDate(pe, pl); got != pe {
			t.Fatalf("EarliestDate: ожидали %v, получили %v", &early, got)
		}
	})
}

func TestParseScheduledAt(t *testing.T) {
	t.Parallel()
	t.Run("пустая_строка_и_пробелы_как_сейчас", func(t *testing.T) {
		for _, in := range []string{"", "   ", "\t"} {
			got, err := ParseScheduledAt(in)
			if err != nil {
				t.Fatalf("ParseScheduledAt: %q: %v", in, err)
			}
			if d := time.Since(got); d < 0 || d > 2*time.Second {
				t.Fatalf("ParseScheduledAt: %q=%v, ожидали около now", in, got)
			}
		}
	})
	t.Run("RFC3339", func(t *testing.T) {
		want := time.Date(2024, 3, 15, 14, 30, 0, 0, time.FixedZone("", 3*3600))
		got, err := ParseScheduledAt("2024-03-15T14:30:00+03:00")
		if err != nil {
			t.Fatalf("ParseScheduledAt: %v", err)
		}
		if !got.Equal(want) {
			t.Fatalf("ParseScheduledAt: ожидали %v, получили %v", want, got)
		}
	})
	t.Run("RFC3339Nano_с_долей_секунды", func(t *testing.T) {
		const s = "2024-03-15T14:30:00.5+03:00"
		want, werr := time.Parse(time.RFC3339Nano, s)
		if werr != nil {
			t.Fatal(werr)
		}
		got, err := ParseScheduledAt(s)
		if err != nil {
			t.Fatalf("ParseScheduledAt: %v", err)
		}
		if !got.Equal(want) {
			t.Fatalf("ParseScheduledAt: ожидали %v, получили %v", want, got)
		}
	})
	t.Run("legacy_дата_время_в_UTC", func(t *testing.T) {
		want := time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC)
		got, err := ParseScheduledAt("2020-01-02 12:00:00")
		if err != nil {
			t.Fatalf("ParseScheduledAt: %v", err)
		}
		if !got.Equal(want) {
			t.Fatalf("ParseScheduledAt: ожидали %v, получили %v", want, got)
		}
	})
	t.Run("некорректная_строка", func(t *testing.T) {
		_, err := ParseScheduledAt("не-дата")
		if !errors.Is(err, ErrWrongScheduledAtFormat) {
			t.Fatalf("ParseScheduledAt: ожидали ErrWrongScheduledAtFormat, получили %v", err)
		}
	})
}

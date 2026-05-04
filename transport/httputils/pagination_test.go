package httputils

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const exampleURL = "http://example.ru/res"

func TestParseListPagination(t *testing.T) {
	t.Parallel()

	t.Run("limit_и_offset_заданы", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		limitExp := 20
		offsetExp := 40
		q := raw.Query()
		q.Set("limit", strconv.Itoa(limitExp))
		q.Set("offset", strconv.Itoa(offsetExp))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != limitExp {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", limitExp, limit)
		}
		if offset != offsetExp {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", offsetExp, offset)
		}
	})
	t.Run("offset=0", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		limitExp := 20
		offsetExp := 0
		q := raw.Query()
		q.Set("limit", strconv.Itoa(limitExp))
		q.Set("offset", strconv.Itoa(offsetExp))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != limitExp {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", limitExp, limit)
		}
		if offset != offsetExp {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", offsetExp, offset)
		}
	})

	t.Run("MaxLimit", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		limitExp := MaxLimit
		q := raw.Query()
		q.Set("limit", strconv.Itoa(MaxLimit))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != limitExp {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", limitExp, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("если_больше_MaxLimit_то_MaxLimit", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		limitExp := MaxLimit
		q := raw.Query()
		q.Set("limit", strconv.Itoa(MaxLimit+1))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != limitExp {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", limitExp, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("дефолт_если_limit_и_offset_не_переданы", func(t *testing.T) {

		req := httptest.NewRequest("GET", exampleURL, nil)
		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != DefaultLimit {
			t.Fatalf("ParseListPagination: limit ожидали дефолт %d, получили%d", DefaultLimit, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})

	t.Run("дефолт_limit_и_корректный_offset", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		offsetExp := 20
		q := raw.Query()

		q.Set("offset", strconv.Itoa(offsetExp))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != DefaultLimit {
			t.Fatalf("ParseListPagination: limit ожидали дефолт %d, получили%d", DefaultLimit, limit)
		}
		if offset != offsetExp {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("некорректный_limit=0", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}

		q := raw.Query()
		q.Set("limit", strconv.Itoa(0))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if !errors.Is(err, ErrInvalidLimit) {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != 0 {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", 0, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("некорректный_limit_меньше_0", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}

		q := raw.Query()
		q.Set("limit", strconv.Itoa(-10))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if !errors.Is(err, ErrInvalidLimit) {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != 0 {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", 0, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("некорректный_offset_меньше_0", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}

		q := raw.Query()
		q.Set("offset", strconv.Itoa(-10))
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if !errors.Is(err, ErrInvalidOffset) {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != 0 {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", 0, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})

	t.Run("некорректный_строковый_limit", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}

		q := raw.Query()
		q.Set("limit", "aa")
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if !errors.Is(err, ErrInvalidLimit) {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != 0 {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", 0, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})
	t.Run("некорректный_строковый_offset", func(t *testing.T) {
		raw, err := url.Parse(exampleURL)
		if err != nil {
			t.Fatalf("ParseListPagination: %v", err)
		}

		q := raw.Query()
		q.Set("offset", "aa")
		raw.RawQuery = q.Encode()
		req := httptest.NewRequest("GET", raw.String(), nil)

		limit, offset, err := ParseListPagination(req)
		if !errors.Is(err, ErrInvalidOffset) {
			t.Fatalf("ParseListPagination: %v", err)
		}
		if limit != 0 {
			t.Fatalf("ParseListPagination: limit ожидали %d, получили%d", 0, limit)
		}
		if offset != 0 {
			t.Fatalf("ParseListPagination: offset ожидали %d, получили%d", 0, offset)
		}
	})

}

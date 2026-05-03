package httputils

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const exampleURL = "http://example.ru/res"

func TestParseListPagination(t *testing.T) {

	t.Run("limit и offset заданы", func(t *testing.T) {

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
			t.Fatalf("ParseListPagination: ожидали %d, получили%d", limitExp, limit)
		}
		if offset != offsetExp {
			t.Fatalf("ParseListPagination: ожидали %d, получили%d", offsetExp, offset)

		}
	})
}

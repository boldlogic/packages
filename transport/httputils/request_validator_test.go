package httputils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type sampleJSON struct {
	Name string `json:"name,omitempty"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var sample sampleJSON
	var err error
	sample, err = DecodeRequest[sampleJSON](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sample.Name))

}

func Test_checkContentType(t *testing.T) {
	p, _ := json.Marshal(sampleJSON{
		Name: "имя",
	})
	r := bytes.NewBuffer(p)

	t.Parallel()

	t.Run("Content-Type:application/json", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/", r)
		req.Header.Add("Content-Type", "application/json")
		err := checkContentType(req)
		if err != nil {
			t.Fatalf("checkContentType: %v", err)
		}
	})
	t.Run("Составной_Content-Type:application/json;charset=utf-8", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/", r)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		err := checkContentType(req)
		if err != nil {
			t.Fatalf("checkContentType: %v", err)
		}
	})
	t.Run("Неподдерживаемый_Content-Type:application/xml", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/", r)
		req.Header.Add("Content-Type", "application/xml")
		err := checkContentType(req)
		if err == nil {
			t.Fatalf("checkContentType: %v", err)
		}
	})
	t.Run("Нет_Content-Type", func(t *testing.T) {
		p, _ := json.Marshal(sampleJSON{
			Name: "имя",
		})
		r := bytes.NewBuffer(p)
		req := httptest.NewRequest("POST", "/", r)
		err := checkContentType(req)
		if err == nil {
			t.Fatalf("checkContentType: %v", err)
		}
	})
}

func Test_readLimitedRequestBody(t *testing.T) {
	t.Run("Есть_тело_в_json", func(t *testing.T) {
		p, _ := json.Marshal(sampleJSON{
			Name: "имя",
		})
		r := bytes.NewBuffer(p)
		req := httptest.NewRequest("POST", "/", r)
		body, err := readLimitedRequestBody(req)
		if err != nil {
			t.Fatalf("readLimitedRequestBody: %v", err)
		}
		if body == nil {
			t.Fatalf("readLimitedRequestBody: body nil, ожидалось не пустое")
		}

	})

	t.Run("Пустое_тело_в_json", func(t *testing.T) {
		r := bytes.NewBufferString("")
		req := httptest.NewRequest("POST", "/", r)
		_, err := readLimitedRequestBody(req)
		if err != nil {
			t.Fatalf("readLimitedRequestBody: %v", err)
		}

	})

}

func TestDecodeRequest(t *testing.T) {
	t.Run("успешный_разбор_и_валидация", func(t *testing.T) {
		payload, err := json.Marshal(sampleJSON{Name: "имя"})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		got, err := DecodeRequest[sampleJSON](req)
		if err != nil {
			t.Fatalf("DecodeRequest: %v", err)
		}
		if got.Name != "имя" {
			t.Fatalf("Name=%q, ожидали имя", got.Name)
		}
	})
}

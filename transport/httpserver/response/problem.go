package response

import "net/http"

// ProblemRFC описывает HTTP-ошибку в формате, близком к RFC 7807.
type ProblemRFC struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// Problem создаёт описание HTTP-ошибки.
// Пустой title заменяется значением по умолчанию для кода статуса.
// Детали внутренних ошибок сервера всегда скрываются.
func Problem(status int, title, detail string) ProblemRFC {
	var p = ProblemRFC{
		Status: status,
		Title:  title,
		Detail: detail,
	}

	if title == "" {
		p.defaultTitle()
	}
	if status == http.StatusInternalServerError {
		p.Detail = "что-то пошло не так"
	}

	return p
}

func (p *ProblemRFC) defaultTitle() {
	switch p.Status {
	case http.StatusNotFound:
		p.Title = "NOT_FOUND"
	case http.StatusBadRequest:
		p.Title = "VALIDATION_ERROR"
	case http.StatusConflict:
		p.Title = "CONFLICT"
	case http.StatusInternalServerError:
		p.Title = "SERVER_ERROR"
	case http.StatusUnprocessableEntity:
		p.Title = "BUSINESS_VALIDATION_ERROR"
	case http.StatusUnsupportedMediaType:
		p.Title = "UNSUPPORTED_MEDIA_TYPE"
	case http.StatusRequestEntityTooLarge:
		p.Title = "REQUEST_ENTITY_TOO_LARGE"
	default:
		p.Title = http.StatusText(p.Status)
	}
}

package shutdown

import (
	"context"
	"errors"
)

// IsExceeded возвращает true, если ошибка означает отмену контекста
// или превышение дедлайна.
//
// Функция корректно работает как с исходными ошибками `context.Canceled`
// и `context.DeadlineExceeded`, так и с ошибками, обёрнутыми через `%w`.
func IsExceeded(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

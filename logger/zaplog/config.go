package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Config содержит параметры конфигурации для создания логгера.
// Level — уровень логирования (debug, info, warn, error, panic).
// Format — формат вывода: "JSON" для JSON-формата, любое другое значение для консольного формата.
// OutputFile — путь к файлу для записи логов. Пустое значение означает вывод в stdout.
type Config struct {
	Level      string `yaml:"level" json:"level"`
	Format     string `yaml:"format" json:"format"`
	OutputFile string `yaml:"output_file" json:"output_file"`
}

// Validate проверяет корректность конфигурации логгера.
// Возвращает список ошибок валидации. Пустой список означает успешную валидацию.
func (c Config) Validate() []error {
	var errs []error
	_, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		errs = append(errs, fmt.Errorf("некорректный уровень логирования %q", c.Level))
	}
	return errs
}

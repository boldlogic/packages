# commonconfig

Пакет `commonconfig` предоставляет утилиты для работы с конфигурационными файлами: получение пути к файлу и декодирование в типизированные структуры.

## Публичный API

```go
func GetConfigPath(defaultConfigPath string) string
func DecodeConfig[T any](path string) (T, error)
```

## Функции

### `GetConfigPath`

Возвращает путь к файлу конфигурации из флага `-config`. Если флаг не передан, возвращает `defaultConfigPath`.

```go
path := GetConfigPath("config/default.yaml")
```

Запуск с кастомным путём:

```bash
go run main.go -config custom.yaml
```

### `DecodeConfig`

Декодирует файл конфигурации в структуру типа `T`. Поддерживает форматы: `.yaml`, `.yml`, `.json`. Выбор формата определяется по расширению файла.

```go
type Config struct {
    Name string `json:"name" yaml:"name"`
    Port int    `json:"port" yaml:"port"`
}

cfg, err := DecodeConfig[Config]("config.yaml")
```

## Ошибки

- `ErrFileStat` — не удалось получить информацию о файле
- `ErrFileRead` — не удалось прочитать файл
- `ErrWrongFileExt` — некорректное расширение файла конфигурации

## Установка

```bash
go get github.com/boldlogic/packages/commonconfig
```

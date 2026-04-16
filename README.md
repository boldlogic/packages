# packages

[![CI](https://github.com/boldlogic/packages/actions/workflows/go.yml/badge.svg)](https://github.com/boldlogic/packages/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/badge/go-1.26.2-blue.svg)](https://golang.org)

Набор небольших Go-пакетов для типовых задач в сервисах и CLI-приложениях: загрузка конфигурации, декодирование JSON/YAML и инициализация логгера на базе `zap`.

## Что есть в репозитории

| Пакет | Назначение |
| --- | --- |
| [commonconfig](./commonconfig) | Получение пути к конфигу и декодирование файла в типизированную структуру |
| [logger/zaplog](./logger/zaplog) | Минимальная обёртка над `zap` с простым конфигом |
| [utils/converters](./utils/converters) | Дженерик-декодеры JSON и YAML из `[]byte` |

## Когда это полезно

- Нужен единый способ читать конфигурацию из `.yaml`, `.yml` или `.json`.
- Хочется быстро поднять структурированный логгер без отдельного слоя инициализации.
- Нужны компактные переиспользуемые пакеты без тяжёлой инфраструктуры.

## Требования

- Go `1.26.2+`

## Установка

```bash
go get github.com/boldlogic/packages@latest
```

Если нужен только один пакет, можно подключать его напрямую:

```bash
go get github.com/boldlogic/packages/commonconfig
go get github.com/boldlogic/packages/logger/zaplog
go get github.com/boldlogic/packages/utils/converters
```

## Быстрый старт

```go
package main

import (
	"log"

	"github.com/boldlogic/packages/commonconfig"
	zaplog "github.com/boldlogic/packages/logger/zaplog"
)

type Config struct {
	Log zaplog.Config `json:"log" yaml:"log"`
}

func main() {
	path := commonconfig.GetConfigPath("config.yaml")

	cfg, err := commonconfig.DecodeConfig[Config](path)
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	if errs := cfg.Log.Validate(); len(errs) > 0 {
		log.Fatalf("некорректная конфигурация логгера: %v", errs)
	}

	logger := zaplog.New(cfg.Log)
	defer logger.Sync()

	logger.Info("сервис запущен")
}
```

Запуск с кастомным конфигом:

```bash
go run ./cmd/app -config ./configs/dev.yaml
```

## Что важно знать

- `commonconfig.DecodeConfig` сохраняет текущее мягкое поведение и игнорирует неизвестные поля.
- Если нужна строгая проверка структуры конфига, используйте `commonconfig.DecodeConfigStrict`.
- `zaplog` пишет либо в `stdout`, либо в один файл, указанный в `OutputFile`.
- Если файл логов открыть не удалось, `zaplog` автоматически переключается на `stdout`.

## Зависимости

- `go.uber.org/zap` для логирования
- `gopkg.in/yaml.v3` для YAML

## Разработка

Основные проверки в репозитории:

```bash
go test ./...
go vet ./...
```

CI в GitHub Actions выполняет сборку, тесты и базовые проверки качества.

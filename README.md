# packages

[![CI](https://github.com/boldlogic/packages/actions/workflows/go.yml/badge.svg)](https://github.com/boldlogic/packages/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/badge/go-1.26.2-blue.svg)](https://golang.org)

Набор небольших Go-пакетов для типовых задач в сервисах и CLI-приложениях: загрузка конфигурации, декодирование JSON/YAML, in-memory cache, базовый Prometheus registry, периодические фоновые задачи, инициализация логгера на базе `zap`, подключение к Microsoft SQL Server и обработка ошибок graceful shutdown.

## Что есть в репозитории

| Пакет | Назначение |
| --- | --- |
| [cache](./cache) | Простой in-memory cache с TTL для записей |
| [commonconfig](./commonconfig) | Получение пути к конфигу и декодирование файла в типизированную структуру |
| [dbzap](./dbzap) | Конфиг подключения к SQL Server и открытие соединения через `database/sql` |
| [logger/zaplog](./logger/zaplog) | Минимальная обёртка над `zap` с простым конфигом |
| [metrics](./metrics) | Базовый Prometheus registry со стандартными коллекторами |
| [periodic](./periodic) | Запуск периодических worker и координация фоновых задач |
| [shutdown](./shutdown) | Распознавание ошибок отмены контекста и превышения дедлайна |
| [utils/converters](./utils/converters) | Дженерик-декодеры JSON и YAML из `[]byte` |
| [utils/xmlconv](./utils/xmlconv) | Утилиты для декодирования XML, включая числа с десятичной запятой |

## Когда это полезно

- Нужен единый способ читать конфигурацию из `.yaml`, `.yml` или `.json`.
- Нужен небольшой локальный cache в памяти с TTL.
- Нужен готовый Prometheus registry без дублирования инициализации стандартных коллекторов.
- Нужно запускать фоновые задачи по интервалу и останавливать их по контексту.
- Хочется быстро поднять структурированный логгер без отдельного слоя инициализации.
- Нужно стандартно описывать подключение к SQL Server и открывать его с `PingContext`.
- Нужно отличать ожидаемую отмену по контексту от настоящих ошибок приложения.
- Нужно разбирать XML с числами в формате `12,34`.
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
go get github.com/boldlogic/packages/cache
go get github.com/boldlogic/packages/dbzap
go get github.com/boldlogic/packages/logger/zaplog
go get github.com/boldlogic/packages/metrics
go get github.com/boldlogic/packages/periodic
go get github.com/boldlogic/packages/shutdown
go get github.com/boldlogic/packages/utils/converters
go get github.com/boldlogic/packages/utils/xmlconv
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
- `cache` не запускает фоновую очистку автоматически: для удаления истёкших записей используется `Cleanup`.
- Если нужна строгая проверка структуры конфига, используйте `commonconfig.DecodeConfigStrict`.
- `dbzap` сейчас ориентирован на `sqlserver` и проверяет соединение с БД через `PingContext` при создании.
- `metrics.New` сразу регистрирует стандартные Go- и process-метрики Prometheus.
- `periodic` запускает worker до отмены контекста и подходит для фоновых сервисных задач.
- `shutdown.IsExceeded` возвращает `true` для `context.Canceled` и `context.DeadlineExceeded`, включая обёрнутые ошибки.
- `zaplog` пишет либо в `stdout`, либо в один файл, указанный в `OutputFile`.
- Если файл логов открыть не удалось, `zaplog` автоматически переключается на `stdout`.
- `xmlconv.RuFloat` помогает читать XML-числа с десятичной запятой без ручного пост-обработчика.

## Зависимости

- `go.uber.org/zap` для логирования
- `github.com/microsoft/go-mssqldb` для SQL Server
- `gopkg.in/yaml.v3` для YAML

## Разработка

Основные проверки в репозитории:

```bash
golangci-lint run
go run ./cmd/projectlint
go test ./...
go vet ./...
```

CI в GitHub Actions выполняет проверку уязвимостей, линтинг, project-specific проверки, сборку и тесты.

# packages

[![CI](https://github.com/boldlogic/packages/actions/workflows/go.yml/badge.svg)](https://github.com/boldlogic/packages/actions/workflows/go.yml)

Набор небольших Go-пакетов для повторно используемых задач: конфигурация, JSON/YAML, cache, Prometheus, HTTP middleware, HTTP-валидация, даты, фоновые задачи, zap-логгер, SQL Server и graceful shutdown.


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
| [transport/httpserver/httpmetrics](./transport/httpserver/httpmetrics) | Prometheus-метрики для HTTP-запросов |
| [transport/httpserver/middleware](./transport/httpserver/middleware) | Middleware для логирования HTTP-запросов и записи метрик |
| [transport/httputils](./transport/httputils) | Разбор JSON-запросов, валидация структур и пагинация из query-параметров |
| [utils/converters](./utils/converters) | Дженерик-декодеры JSON и YAML из `[]byte` |
| [utils/dates](./utils/dates) | Утилиты для разбора, нормализации и сравнения дат |
| [utils/xmlconv](./utils/xmlconv) | Утилиты для декодирования XML, включая числа с десятичной запятой |

## Требования

- Версия Go указана в [go.mod](./go.mod).

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
go get github.com/boldlogic/packages/transport/httpserver/httpmetrics
go get github.com/boldlogic/packages/transport/httpserver/middleware
go get github.com/boldlogic/packages/transport/httputils
go get github.com/boldlogic/packages/utils/converters
go get github.com/boldlogic/packages/utils/dates
go get github.com/boldlogic/packages/utils/xmlconv
```

## Разработка

Основные проверки:

```bash
golangci-lint run
go run ./cmd/projectlint
go test ./...
go vet ./...
```

CI выполняет `golangci-lint`, `projectlint`, сборку и тесты.

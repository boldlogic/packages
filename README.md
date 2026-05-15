# packages

[CI](https://github.com/boldlogic/packages/actions/workflows/go.yml)

Набор небольших Go-пакетов для повторно используемых задач: конфигурация, JSON/YAML, cache, Prometheus, HTTP middleware, HTTP-ответы, HTTP-валидация, даты, фоновые задачи, zap-логгер, SQL Server и graceful shutdown.

## Что есть в репозитории


| Пакет                                                                  | Назначение                                                                 |
| ---------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| [cache](./cache)                                                       | Простой in-memory cache с TTL для записей                                  |
| [commonconfig](./commonconfig)                                         | Получение пути к конфигу и декодирование файла в типизированную структуру  |
| [dbzap](./dbzap)                                                       | Конфиг подключения к SQL Server и открытие соединения через `database/sql` |
| [logger/zaplog](./logger/zaplog)                                       | Минимальная обёртка над `zap` с простым конфигом                           |
| [metrics](./metrics)                                                   | Базовый Prometheus registry со стандартными коллекторами                   |
| [periodic](./periodic)                                                 | Запуск периодических worker и координация фоновых задач                    |
| [shutdown](./shutdown)                                                 | Распознавание ошибок отмены контекста и превышения дедлайна                |
| [transport/httpserver](./transport/httpserver)                         | Обёртка над `net/http.Server` с типизированной конфигурацией и дефолтами    |
| [transport/httpserver/httpmetrics](./transport/httpserver/httpmetrics) | Prometheus-метрики для HTTP-запросов                                       |
| [transport/httpserver/middleware](./transport/httpserver/middleware)   | Middleware для логирования HTTP-запросов и записи метрик                   |
| [transport/httpserver/response](./transport/httpserver/response)       | Запись JSON-ответов HTTP и описание ошибок в формате, близком к RFC 7807   |
| [transport/httputils](./transport/httputils)                           | Разбор JSON-запросов, валидация структур и пагинация из query-параметров   |
| [utils/converters](./utils/converters)                                 | Дженерик-декодеры JSON и YAML из `[]byte`                                  |
| [utils/dates](./utils/dates)                                           | Утилиты для разбора, нормализации и сравнения дат                          |
| [utils/xmlconv](./utils/xmlconv)                                       | Утилиты для декодирования XML, включая числа с десятичной запятой          |


## Установка

```bash
go get github.com/boldlogic/packages@<версия>
```

Импортируйте нужные подпакеты по путям из таблицы выше.

## Разработка

Основные проверки:

```bash
golangci-lint run
go run ./cmd/projectlint
go test ./...
go vet ./...
```

CI выполняет `golangci-lint`, `projectlint`, сборку и тесты.

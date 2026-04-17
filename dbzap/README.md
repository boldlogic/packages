# dbzap

`dbzap` — небольшой пакет для подключения к Microsoft SQL Server через `database/sql` и `github.com/microsoft/go-mssqldb`.

Пакет решает две задачи:
- хранит конфиг подключения в одной структуре;
- открывает соединение с проверкой доступности базы через `PingContext`.

## Публичный API

```go
type DBConfig struct {
    Driver   string
    Server   string
    Host     string
    Port     int
    Name     string
    User     string
    Password string
    SSLMode  string
}

func (db *DBConfig) ApplyDefaults()
func (db *DBConfig) ApplySecretsFromEnv()
func (db *DBConfig) Validate() []error
func (db *DBConfig) GetDSN() string

type Pool struct {
    Db     *sql.DB
    Logger *zap.Logger
}

func New(ctx context.Context, dsn string, logger *zap.Logger) (*Pool, error)
func (p *Pool) Close()
```

## Ошибки

- `ErrDBOpen` — не удалось открыть подключение к БД
- `ErrDBPing` — не удалось проверить подключение к БД

Ошибки из `New` и `openDB` оборачиваются через `%w`, поэтому их можно проверять через `errors.Is`.

## Быстрый старт

```go
package main

import (
	"context"
	"log"

	"github.com/boldlogic/packages/dbzap"
	zaplog "github.com/boldlogic/packages/logger/zaplog"
)

func main() {
	cfg := dbzap.DBConfig{
		Server:   "localhost",
		Name:     "app",
		User:     "sa",
		Password: "secret",
	}

	cfg.ApplyDefaults()
	cfg.ApplySecretsFromEnv()

	if errs := cfg.Validate(); len(errs) > 0 {
		log.Fatalf("некорректная конфигурация БД: %v", errs)
	}

	logger := zaplog.New(zaplog.Config{
		Level:  "info",
		Format: "json",
	})
	defer logger.Sync()

	pool, err := dbzap.New(context.Background(), cfg.GetDSN(), logger)
	if err != nil {
		log.Fatalf("не удалось создать подключение к БД: %v", err)
	}
	defer pool.Close()
}
```

## Переменные окружения

`ApplySecretsFromEnv` поддерживает:

- `DB_PASSWORD`
- `MSSQL_SA_PASSWORD`
- `DB_USER`

Приоритет для пароля:
1. `DB_PASSWORD`
2. `MSSQL_SA_PASSWORD`

## Что важно знать

- Текущий пакет ориентирован на `sqlserver`.
- Для валидации обязательным считается поле `Server`.
- `Host`, `Port` и `SSLMode` сохраняются в конфиге, но в текущем DSN не используются.
- `New` не делает retry и не настраивает параметры пула соединений.

## Установка

```bash
go get github.com/boldlogic/packages/dbzap
```

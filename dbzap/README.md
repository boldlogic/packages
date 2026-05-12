# dbzap

`dbzap` — небольшой пакет для подключения к Microsoft SQL Server через `database/sql` и `github.com/microsoft/go-mssqldb`.

## Ошибки

- `ErrDBOpen` — не удалось открыть подключение к БД
- `ErrDBPing` — не удалось проверить подключение к БД

Ошибки из `New` и `openDB` оборачиваются через `%w`, поэтому их можно проверять через `errors.Is`.

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

## Установка

```bash
go get github.com/boldlogic/packages/dbzap
```

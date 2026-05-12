# shutdown

`shutdown` — небольшой пакет с утилитой для распознавания ошибок завершения по контексту.

## Что делает пакет

- возвращает `true` для `context.Canceled`
- возвращает `true` для `context.DeadlineExceeded`
- корректно работает с обёрнутыми ошибками через `%w`

## Установка

```bash
go get github.com/boldlogic/packages/shutdown
```

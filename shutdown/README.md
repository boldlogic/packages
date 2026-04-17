# shutdown

`shutdown` — небольшой пакет с утилитой для распознавания ошибок завершения по контексту.

Он полезен в коде graceful shutdown, фоновых воркерах и сетевых обработчиках, где отмена контекста и превышение дедлайна считаются ожидаемыми сценариями, а не полноценными ошибками приложения.

## Публичный API

```go
func IsExceeded(err error) bool
```

## Что делает пакет

- возвращает `true` для `context.Canceled`
- возвращает `true` для `context.DeadlineExceeded`
- корректно работает с обёрнутыми ошибками через `%w`

## Пример

```go
if err := worker(ctx); err != nil {
	if shutdown.IsExceeded(err) {
		return
	}
	log.Printf("неожиданная ошибка: %v", err)
}
```

## Установка

```bash
go get github.com/boldlogic/packages/shutdown
```

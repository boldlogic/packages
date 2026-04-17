# cache

`cache` — простой in-memory cache с TTL для записей.

Пакет подходит для небольших локальных кешей в приложении, когда не нужен внешний storage, фоновые воркеры или сложная политика вытеснения.

## Публичный API

```go
type CacheItem[T any] struct {
    Value T
    Ttl   int64
}

type Cache[T any] struct {}

func NewCache[T any](ttl time.Duration) *Cache[T]
func (c *Cache[T]) Get(key string) (T, bool)
func (c *Cache[T]) Set(key string, value T)
func (c *Cache[T]) Evict(key string)
func (c *Cache[T]) Cleanup()
```

## Что делает пакет

- хранит значения в памяти по строковому ключу
- задаёт общий TTL для всех записей при создании cache
- возвращает `false` из `Get`, если ключ отсутствует или запись истекла
- умеет явно удалять запись через `Evict`
- умеет очищать все истёкшие записи через `Cleanup`

## Пример

```go
c := cache.NewCache[string](time.Minute)

c.Set("token", "abc123")

value, ok := c.Get("token")
if !ok {
	return
}

fmt.Println(value)
```

## Что важно знать

- пакет не запускает фоновую очистку автоматически
- истёкшие записи удаляются лениво через `Get` или явно через `Cleanup`
- пакет потокобезопасен для конкурентного доступа

## Установка

```bash
go get github.com/boldlogic/packages/cache
```

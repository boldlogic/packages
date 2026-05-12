# converters

Пакет `converters` предоставляет универсальные функции для декодирования данных в форматах JSON и YAML в типизированные структуры Go.

## Что делает пакет

- декодирует JSON из `[]byte` в тип `T`
- декодирует YAML из `[]byte` в тип `T`
- strict JSON запрещает неизвестные поля и несколько JSON-значений
- strict YAML запрещает неизвестные поля и несколько YAML-документов

## Ошибки

- `ErrWrongJSON` — некорректный формат JSON
- `ErrWrongYAML` — некорректный формат YAML

## Установка

```bash
go get github.com/boldlogic/packages/utils/converters
```

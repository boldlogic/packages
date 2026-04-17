# xmlconv

`xmlconv` — небольшой пакет с утилитами для декодирования XML в Go-структуры.

Сейчас пакет содержит тип `RuFloat`, который помогает читать числа из XML, где дробная часть может быть записана через запятую.

## Публичный API

```go
type RuFloat float64

func (rf *RuFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
```

## Что делает пакет

- поддерживает числа вида `12,34`
- поддерживает числа вида `12.34`
- интерпретирует пустое значение как `0`
- работает через стандартный `encoding/xml`

## Пример

```go
type Item struct {
	Price xmlconv.RuFloat `xml:"Price"`
}
```

Если в XML придёт:

```xml
<Price>12,34</Price>
```

то после декодирования значение будет равно `12.34`.

## Установка

```bash
go get github.com/boldlogic/packages/utils/xmlconv
```

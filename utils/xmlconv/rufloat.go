package xmlconv

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// RuFloat представляет число с поддержкой XML-значений, где дробная часть
// может быть отделена запятой.
type RuFloat float64

// UnmarshalXML декодирует XML-элемент в число с плавающей точкой.
//
// Функция поддерживает значения с запятой и точкой в роли десятичного
// разделителя. Пустое значение интерпретируется как `0`.
func (rf *RuFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		*rf = RuFloat(0)
		return nil
	}
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*rf = RuFloat(f)
	return nil
}

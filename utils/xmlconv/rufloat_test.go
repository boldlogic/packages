package xmlconv

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestRuFloat_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name   string
		xmlStr string
		want   float64
	}{
		{
			name:   "запятая",
			xmlStr: "<Value>12,34</Value>",
			want:   12.34,
		},
		{
			name:   "точка",
			xmlStr: "<Value>12.34</Value>",
			want:   12.34,
		},
		{
			name:   "пустое_значение",
			xmlStr: "<Value></Value>",
			want:   0,
		},
		{
			name:   "целое_число",
			xmlStr: "<Value>100</Value>",
			want:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			full := "<?xml version=\"1.0\"?><root>" + tt.xmlStr + "</root>"
			var wrapper struct {
				Value RuFloat `xml:"Value"`
			}

			dec := xml.NewDecoder(strings.NewReader(full))
			if err := dec.Decode(&wrapper); err != nil {
				t.Fatalf("декодирование XML: %v", err)
			}
			if float64(wrapper.Value) != tt.want {
				t.Fatalf("получили %v, ожидали %v", float64(wrapper.Value), tt.want)
			}
		})
	}
}

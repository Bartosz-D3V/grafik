package common

import (
	"bytes"
	"unicode"
)

func SentenceCase(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func SnakeCaseToCamelCase(s string) string {
	if s == "" {
		return s
	}
	var buff bytes.Buffer
	buff.WriteByte(s[0])
	capitalize := false
	for i := 1; i < len(s); i++ {
		char := s[i]
		if char == '_' {
			capitalize = true
			continue
		} else if capitalize {
			capitalize = false
			buff.WriteRune(unicode.ToUpper(rune(char)))
		} else {
			capitalize = false
			buff.WriteByte(char)
		}
	}
	return buff.String()
}

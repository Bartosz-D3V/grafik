// Package common contains commonly used helper functions used within grafik project.
package common

import (
	"bytes"
	"unicode"
)

// SentenceCase returns a new string with first character lowercased.
// "MyExample" -> "myExample".
// "iPhone" -> "iPhone".
// "1Note" -> "1Note".
func SentenceCase(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// SnakeCaseToCamelCase converts snake_case to CamelCase.
// "my_example" -> "MyExample".
// "i_phone" -> "iPhone".
// "1_Note" -> "1Note".
func SnakeCaseToCamelCase(s string) string {
	if s == "" {
		return s
	}
	var buff bytes.Buffer
	buff.WriteByte(s[0])
	capitalize := false
	for i := 1; i < len(s); i++ {
		char := s[i]
		switch {
		case char == '_':
			capitalize = true
		case capitalize:
			capitalize = false
			buff.WriteRune(unicode.ToUpper(rune(char)))
		default:
			capitalize = false
			buff.WriteByte(char)
		}
	}
	return buff.String()
}

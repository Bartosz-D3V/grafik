package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeArg_ExportName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		t   TypeArg
		exp string
	}{
		{TypeArg{Name: "test"}, "Test"},
		{TypeArg{Name: "Test"}, "Test"},
		{TypeArg{Name: "1test"}, "1test"},
		{TypeArg{Name: "1Test"}, "1Test"},
		{TypeArg{Name: ""}, ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.t.ExportName())
	}
}

func TestTypeArg_ExportType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		t   TypeArg
		exp string
	}{
		{TypeArg{Type: "string"}, "string"},
		{TypeArg{Type: "int"}, "int"},
		{TypeArg{Type: "bool"}, "bool"},
		{TypeArg{Type: "[]string"}, "[]string"},
		{TypeArg{Type: "[]int"}, "[]int"},
		{TypeArg{Type: "[]bool"}, "[]bool"},
		{TypeArg{Type: "[][]string"}, "[][]string"},
		{TypeArg{Type: "[][]int"}, "[][]int"},
		{TypeArg{Type: "[][]bool"}, "[][]bool"},
		{TypeArg{Type: "[][][]string"}, "[][][]string"},
		{TypeArg{Type: "[][][]int"}, "[][][]int"},
		{TypeArg{Type: "[][][]bool"}, "[][][]bool"},
		{TypeArg{Type: "person"}, "Person"},
		{TypeArg{Type: "[]person"}, "[]Person"},
		{TypeArg{Type: "[][]person"}, "[][]Person"},
		{TypeArg{Type: "[][][]person"}, "[][][]Person"},
		{TypeArg{Type: "Person"}, "Person"},
		{TypeArg{Type: "[]Person"}, "[]Person"},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.t.ExportType().Type)
	}
}

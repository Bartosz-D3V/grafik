package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeField_ExportName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		t   TypeField
		exp string
	}{
		{TypeField{Name: "test"}, "Test"},
		{TypeField{Name: "Test"}, "Test"},
		{TypeField{Name: "1test"}, "1test"},
		{TypeField{Name: "1Test"}, "1Test"},
		{TypeField{Name: ""}, ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.t.ExportName())
	}
}

func TestTypeField_ExportType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		t   TypeField
		exp string
	}{
		{TypeField{Type: "string"}, "string"},
		{TypeField{Type: "int"}, "int"},
		{TypeField{Type: "bool"}, "bool"},
		{TypeField{Type: "[]string"}, "[]string"},
		{TypeField{Type: "[]int"}, "[]int"},
		{TypeField{Type: "[]bool"}, "[]bool"},
		{TypeField{Type: "[][]string"}, "[][]string"},
		{TypeField{Type: "[][]int"}, "[][]int"},
		{TypeField{Type: "[][]bool"}, "[][]bool"},
		{TypeField{Type: "[][][]string"}, "[][][]string"},
		{TypeField{Type: "[][][]int"}, "[][][]int"},
		{TypeField{Type: "[][][]bool"}, "[][][]bool"},
		{TypeField{Type: "person"}, "Person"},
		{TypeField{Type: "[]person"}, "[]Person"},
		{TypeField{Type: "[][]person"}, "[][]Person"},
		{TypeField{Type: "[][][]person"}, "[][][]Person"},
		{TypeField{Type: "Person"}, "Person"},
		{TypeField{Type: "[]Person"}, "[]Person"},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.t.ExportType().Type)
	}
}

func TestTypeField_PointerType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		t   TypeField
		exp string
	}{
		{TypeField{Type: "string"}, "*string"},
		{TypeField{Type: "int"}, "*int"},
		{TypeField{Type: "bool"}, "*bool"},
		{TypeField{Type: "[]string"}, "[]string"},
		{TypeField{Type: "[]int"}, "[]int"},
		{TypeField{Type: "[]bool"}, "[]bool"},
		{TypeField{Type: "[][]string"}, "[][]string"},
		{TypeField{Type: "[][]int"}, "[][]int"},
		{TypeField{Type: "[][]bool"}, "[][]bool"},
		{TypeField{Type: "[][][]string"}, "[][][]string"},
		{TypeField{Type: "[][][]int"}, "[][][]int"},
		{TypeField{Type: "[][][]bool"}, "[][][]bool"},
		{TypeField{Type: "person"}, "*person"},
		{TypeField{Type: "[]person"}, "[]person"},
		{TypeField{Type: "[][]person"}, "[][]person"},
		{TypeField{Type: "[][][]person"}, "[][][]person"},
		{TypeField{Type: "Person"}, "*Person"},
		{TypeField{Type: "[]Person"}, "[]Person"},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.t.PointerType().Type)
	}
}

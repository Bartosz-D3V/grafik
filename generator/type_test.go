package generator

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

func TestFunc_JoinArgsBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		f   Func
		exp string
	}{
		{Func{}, ""},
		{
			Func{
				Args: []TypeArg{
					{
						Name: "age",
						Type: "int",
					},
				},
			},
			"age int",
		},
		{
			Func{
				Args: []TypeArg{
					{
						Name: "age",
						Type: "int",
					},
					{
						Name: "name",
						Type: "string",
					},
				},
			},
			"age int, name string",
		},
		{
			Func{
				Args: []TypeArg{
					{
						Name: "age",
						Type: "int",
					},
					{
						Name: "name",
						Type: "string",
					},
					{
						Name: "address",
						Type: "Address",
					},
				},
			},
			"age int, name string, address Address",
		},
		{
			Func{
				Args: []TypeArg{
					{
						Name: "age",
						Type: "int",
					},
					{
						Name: "name",
						Type: "string",
					},
					{
						Name: "address",
						Type: "Address_and_contact_information",
					},
				},
			},
			"age int, name string, address AddressAndContactInformation",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.f.JoinArgsBy(", "))
	}
}

func TestFunc_ExportName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		f   Func
		exp string
	}{
		{Func{Name: "test"}, "Test"},
		{Func{Name: "Test"}, "Test"},
		{Func{Name: "1test"}, "1test"},
		{Func{Name: "1Test"}, "1Test"},
		{Func{Name: ""}, ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.f.ExportName())
	}
}

func TestConst_IsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c   Const
		exp bool
	}{
		{Const{Val: "test"}, true},
		{Const{Val: 1}, false},
		{Const{Val: 1.2}, false},
		{Const{Val: true}, false},
		{Const{Val: false}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.c.IsString())
	}
}

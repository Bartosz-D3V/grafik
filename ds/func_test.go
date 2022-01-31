package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

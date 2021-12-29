package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/ast"
	"testing"
)

func TestIsList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		val *ast.Type
		exp bool
	}{
		{
			val: ast.ListType(
				&ast.Type{
					NamedType: "Person",
				}, &ast.Position{},
			),
			exp: true,
		},
		{
			val: ast.NamedType(
				"Person",
				&ast.Position{},
			),
			exp: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, IsList(test.val), test.exp)
	}
}

func TestIsComplex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		val *ast.Type
		exp bool
	}{
		{
			val: &ast.Type{NamedType: "String"},
			exp: false,
		},
		{
			val: &ast.Type{NamedType: "Int"},
			exp: false,
		},
		{
			val: &ast.Type{NamedType: "ID"},
			exp: false,
		},
		{
			val: &ast.Type{NamedType: "Float"},
			exp: false,
		},
		{
			val: &ast.Type{NamedType: "Boolean"},
			exp: false,
		},
		{
			val: &ast.Type{NamedType: "Person"},
			exp: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, IsComplex(test.val), test.exp)
	}
}

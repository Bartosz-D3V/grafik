package common

import "github.com/vektah/gqlparser/ast"

// IsList determines if ast.Type is a list (single/multi dimensional).
func IsList(astType *ast.Type) bool {
	return astType.IsCompatible(ast.ListType(astType.Elem, astType.Position))
}

// IsComplex determines if GraphQL type is not 'primitive'.
func IsComplex(t *ast.Type) bool {
	return t.NamedType != "String" && t.NamedType != "Int" &&
		t.NamedType != "ID" && t.NamedType != "Float" &&
		t.NamedType != "Boolean"
}

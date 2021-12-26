// Package common contains commonly used helper functions used within grafik project.

package common

import "github.com/vektah/gqlparser/ast"

// IsList determines if ast.Type is a list (single/multi dimensional).
func IsList(astType *ast.Type) bool {
	return astType.IsCompatible(ast.ListType(astType.Elem, astType.Position))
}

// IsComplex determines if GraphQL type is not 'primitive'.
func IsComplex(astType *ast.Type) bool {
	return astType.NamedType != "String" && astType.NamedType != "Int" &&
		astType.NamedType != "ID" && astType.NamedType != "Float" &&
		astType.NamedType != "Boolean"
}

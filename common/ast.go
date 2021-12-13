package common

import "github.com/vektah/gqlparser/ast"

// IsList determines if ast.Type is a list (single/multi dimensional)
func IsList(astType *ast.Type) bool {
	return astType.IsCompatible(ast.ListType(astType.Elem, astType.Position))
}

// FilterCustomFields returns slice of FieldDefinition created by user in .graphql files
func FilterCustomFields(fieldList ast.FieldList) []ast.FieldDefinition {
	cFields := make([]ast.FieldDefinition, 0)
	for _, field := range fieldList {
		if field.Position != nil && field.Position.Src != nil && !field.Position.Src.BuiltIn {
			cFields = append(cFields, *field)
		}
	}
	return cFields
}

// IsComplex determines if GraphQL type is not 'primitive'.
func IsComplex(t *ast.Type) bool {
	return t.NamedType != "String" && t.NamedType != "Int" &&
		t.NamedType != "ID" && t.NamedType != "Float" &&
		t.NamedType != "Boolean"
}

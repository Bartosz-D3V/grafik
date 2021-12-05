package common

import "github.com/vektah/gqlparser/ast"

func IsList(astType *ast.Type) bool {
	return astType.IsCompatible(ast.ListType(astType.Elem, astType.Position))
}

func FilterCustomFields(fieldList ast.FieldList) []ast.FieldDefinition {
	cFields := make([]ast.FieldDefinition, 0)
	for _, field := range fieldList {
		if field.Position != nil && field.Position.Src != nil && !field.Position.Src.BuiltIn {
			cFields = append(cFields, *field)
		}
	}
	return cFields
}

func IsComplex(t *ast.Type) bool {
	return t.NamedType != "String" && t.NamedType != "Int" &&
		t.NamedType != "ID" && t.NamedType != "Float" &&
		t.NamedType != "Boolean"
}

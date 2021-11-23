package common

import "github.com/vektah/gqlparser/ast"

func IsList(astType *ast.Type) bool {
	if astType.IsCompatible(ast.ListType(astType.Elem, astType.Position)) {
		return true
	}
	return false
}

func NumOfBuiltIns(query *ast.Definition) int {
	if query.OneOf("Query") {
		return len(query.Fields) - 2
	} else {
		return len(query.Fields)
	}
}

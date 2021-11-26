package common

import "github.com/vektah/gqlparser/ast"

func IsList(astType *ast.Type) bool {
	return astType.IsCompatible(ast.ListType(astType.Elem, astType.Position))
}

func NumOfBuiltIns(query *ast.Definition) int {
	if query.OneOf("Query") {
		const numOfPredefinedQueryDefs = 2
		return len(query.Fields) - numOfPredefinedQueryDefs
	} else {
		return len(query.Fields)
	}
}

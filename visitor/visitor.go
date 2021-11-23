package visitor

import (
	"github.com/vektah/gqlparser/ast"
)

type Visitor interface {
	IntrospectTypes() []string
}

type visitor struct {
	schema        *ast.Schema
	queryDocument *ast.QueryDocument
	customTypes   []string
}

func New(schema *ast.Schema, queryDocument *ast.QueryDocument) Visitor {
	return &visitor{
		schema:        schema,
		queryDocument: queryDocument,
		customTypes:   make([]string, 0),
	}
}

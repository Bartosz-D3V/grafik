// Package visitor abstracts logic responsible for determining which custom types from GraphQL Schema file should be generated based on usage in GraphQL query file.
package visitor

import (
	"github.com/vektah/gqlparser/ast"
)

// A Visitor is an interface that provides contract for visitor struct and is being used instead of a pointer.
type Visitor interface {
	IntrospectTypes() map[string][]string
}

// visitor is a private struct that can be created with New function.
type visitor struct {
	schema        *ast.Schema
	queryDocument *ast.QueryDocument
	customTypes   map[string][]string
}

// New return instance of visitor - it requires parsed GraphQL schema and query files.
func New(schema *ast.Schema, queryDocument *ast.QueryDocument) Visitor {
	return &visitor{
		schema:        schema,
		queryDocument: queryDocument,
		customTypes:   make(map[string][]string),
	}
}

func (v *visitor) IntrospectTypes() map[string][]string {
	if v.queryDocument.Operations != nil {
		v.parseOpTypes(v.queryDocument.Operations)
	}

	return v.customTypes
}

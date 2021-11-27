package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/Bartosz-D3V/ggrafik/visitor"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

type Evaluator interface {
	Generate() []byte
}

type evaluator struct {
	generator     generator.Generator
	visitor       visitor.Visitor
	schema        *ast.Schema
	queryDocument *ast.QueryDocument
	clientName    string
	queryFuncs    []generator.Func
}

func New(fptr string, schema *ast.Schema, queryDocument *ast.QueryDocument, clientName string) Evaluator {
	return &evaluator{
		generator:     generator.New(fptr),
		visitor:       visitor.New(schema, queryDocument),
		schema:        schema,
		queryDocument: queryDocument,
		clientName:    parseClientName(clientName, schema),
		queryFuncs:    make([]generator.Func, 0),
	}
}

func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()

	e.generator.WriteLineBreak(2)

	e.genSchemaDef()
	e.generator.WriteLineBreak(1)

	e.genOperations()
	e.generator.WriteLineBreak(1)

	e.genClientCode()
	e.generator.WriteLineBreak(1)

	return e.generator.Generate()
}

func parseClientName(name string, schema *ast.Schema) string {
	if name != "" {
		return name
	}
	return strings.Split(schema.Query.Position.Src.Name, ".")[0]
}

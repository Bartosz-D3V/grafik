package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/Bartosz-D3V/ggrafik/visitor"
	"github.com/vektah/gqlparser/ast"
)

type Evaluator interface {
	Generate() []byte
}

type evaluator struct {
	generator     generator.Generator
	visitor       visitor.Visitor
	schema        *ast.Schema
	queryDocument *ast.QueryDocument
}

func New(fptr string, schema *ast.Schema, queryDocument *ast.QueryDocument) Evaluator {
	return &evaluator{
		generator:     generator.New(fptr),
		visitor:       visitor.New(schema, queryDocument),
		schema:        schema,
		queryDocument: queryDocument,
	}
}

func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()

	e.generator.WriteLineBreak(2)

	e.genSchemaDef()
	e.generator.WriteLineBreak(1)

	e.genOperations()

	return e.generator.Generate()
}

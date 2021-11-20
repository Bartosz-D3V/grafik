package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/vektah/gqlparser/ast"
)

type Evaluator interface {
	Generate() []byte
}

type evaluator struct {
	generator     generator.Generator
	schema        *ast.Schema
	queryDocument *ast.QueryDocument
	customTypes   map[string]bool
}

func New(fptr string, schema *ast.Schema, queryDocument *ast.QueryDocument) Evaluator {
	return &evaluator{
		generator:     generator.New(fptr),
		schema:        schema,
		queryDocument: queryDocument,
		customTypes:   make(map[string]bool, 0),
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

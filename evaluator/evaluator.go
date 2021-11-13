package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/vektah/gqlparser/ast"
)

type Evaluator interface {
	Generate() []byte
}

type evaluator struct {
	generator *generator.Generator
	schema    *ast.Schema
	fptr      string
}

func New(fptr string, schema *ast.Schema) Evaluator {
	return &evaluator{
		generator: generator.New(fptr),
		schema:    schema,
	}
}

func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()

	e.generator.WriteLineBreaks(2)

	e.genSchemaDef()
	e.generator.WriteLineBreaks(1)
	return e.generator.Generate()
}

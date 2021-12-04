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
	clientName    string
	pkgName       string
}

func New(rootLoc string, schema *ast.Schema, queryDocument *ast.QueryDocument, clientName string, pkgName string) Evaluator {
	return &evaluator{
		generator:     generator.New(rootLoc),
		visitor:       visitor.New(schema, queryDocument),
		schema:        schema,
		queryDocument: queryDocument,
		clientName:    clientName,
		pkgName:       pkgName,
	}
}

func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()
	e.generator.WriteLineBreak(2)

	e.generator.WritePackage(e.pkgName)
	e.generator.WriteLineBreak(2)

	e.generator.WriteImports()
	e.generator.WriteLineBreak(2)

	e.genSchemaDef()
	e.generator.WriteLineBreak(1)

	e.genOperations()
	e.generator.WriteLineBreak(1)

	e.genClientCode()
	e.generator.WriteLineBreak(1)

	return e.generator.Generate()
}

// Package evaluator contains the logic responsible for evaluating schema & query GraphQL Abstract Syntax Tree [AST]
// It orchestrates the generation of the code by using visitor & generator packages.
package evaluator

import (
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/Bartosz-D3V/grafik/visitor"
	"github.com/vektah/gqlparser/ast"
)

// Evaluator provides a contract for evaluator and is being used as a return type of New instead of a pointer
type Evaluator interface {
	Generate() []byte
}

type evaluator struct {
	generator      generator.Generator
	visitor        visitor.Visitor
	schema         *ast.Schema
	queryDocument  *ast.QueryDocument
	AdditionalInfo AdditionalInfo
}

func New(rootLoc string, schema *ast.Schema, queryDocument *ast.QueryDocument, additionalInfo AdditionalInfo) Evaluator {
	return &evaluator{
		generator:      generator.New(rootLoc),
		visitor:        visitor.New(schema, queryDocument),
		schema:         schema,
		queryDocument:  queryDocument,
		AdditionalInfo: additionalInfo,
	}
}

func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()
	e.generator.WriteLineBreak(2)

	e.generator.WritePackage(e.AdditionalInfo.PackageName)
	e.generator.WriteLineBreak(2)

	e.generator.WriteImports()
	e.generator.WriteLineBreak(2)

	e.genSchemaDef(e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(1)

	e.genOperations()
	e.generator.WriteLineBreak(1)

	e.genClientCode()
	e.generator.WriteLineBreak(1)

	return e.generator.Generate()
}

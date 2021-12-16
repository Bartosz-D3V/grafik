// Package evaluator contains the logic responsible for evaluating schema & query GraphQL Abstract Syntax Tree [AST].
// It orchestrates the generation of the code by using visitor & generator packages.
package evaluator

import (
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/Bartosz-D3V/grafik/visitor"
	"github.com/vektah/gqlparser/ast"
)

// Evaluator provides a contract for evaluator and is being used as a return type of New instead of a pointer.
type Evaluator interface {
	Generate() []byte
}

// evaluator is a struct used internally by grafikgen to wrap all required services and properties.
type evaluator struct {
	generator      generator.Generator // Instance of generator to abstract logic responsible for generating Go code.
	visitor        visitor.Visitor     // Instance of visitor to obtain list of all used custom GraphQL types.
	schema         *ast.Schema         // GraphQL schema provided via CLI.
	queryDocument  *ast.QueryDocument  // GraphQL query document provided via CLI.
	AdditionalInfo AdditionalInfo      // Additional info provided via CLI.
}

// New function creates an instance of evaluator.
func New(rootLoc string, schema *ast.Schema, queryDocument *ast.QueryDocument, additionalInfo AdditionalInfo) Evaluator {
	return &evaluator{
		generator:      generator.New(rootLoc),
		visitor:        visitor.New(schema, queryDocument),
		schema:         schema,
		queryDocument:  queryDocument,
		AdditionalInfo: additionalInfo,
	}
}

// Generate is a root level function that generates the whole grafik client.
func (e *evaluator) Generate() []byte {
	e.generator.WriteHeader()
	e.generator.WriteLineBreak(twoLinesBreak)

	e.generator.WritePackage(e.AdditionalInfo.PackageName)
	e.generator.WriteLineBreak(twoLinesBreak)

	e.generator.WriteImports()
	e.generator.WriteLineBreak(twoLinesBreak)

	e.genSchemaDef(e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(oneLineBreak)

	e.genOperations()
	e.generator.WriteLineBreak(oneLineBreak)

	e.genClientCode()
	e.generator.WriteLineBreak(oneLineBreak)

	return e.generator.Generate()
}

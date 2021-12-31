// Package evaluator contains the logic responsible for evaluating schema & query GraphQL Abstract Syntax Tree [AST].
// It orchestrates the generation of the code by using visitor & generator packages.
package evaluator

import (
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/Bartosz-D3V/grafik/visitor"
	"github.com/vektah/gqlparser/ast"
	"io"
)

// Evaluator provides a contract for evaluator and is being used as a return type of New instead of a pointer.
type Evaluator interface {
	Generate() (io.WriterTo, error)
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
func New(schema *ast.Schema, queryDocument *ast.QueryDocument, additionalInfo AdditionalInfo) (Evaluator, error) {
	g, err := generator.New()
	if err != nil {
		return nil, err
	}
	return &evaluator{
		generator:      g,
		visitor:        visitor.New(schema, queryDocument),
		schema:         schema,
		queryDocument:  queryDocument,
		AdditionalInfo: additionalInfo,
	}, nil
}

// Generate is a root level function that generates the whole grafik client.
func (e *evaluator) Generate() (io.WriterTo, error) {
	err := e.generator.WriteHeader()
	if err != nil {
		return nil, err
	}
	err = e.generator.WriteLineBreak(twoLinesBreak)
	if err != nil {
		return nil, err
	}

	err = e.generator.WritePackage(e.AdditionalInfo.PackageName)
	if err != nil {
		return nil, err
	}
	err = e.generator.WriteLineBreak(twoLinesBreak)
	if err != nil {
		return nil, err
	}

	err = e.generator.WriteImports()
	if err != nil {
		return nil, err
	}

	err = e.generator.WriteLineBreak(twoLinesBreak)
	if err != nil {
		return nil, err
	}

	err = e.genSchemaDef()
	if err != nil {
		return nil, err
	}
	err = e.generator.WriteLineBreak(oneLineBreak)
	if err != nil {
		return nil, err
	}

	err = e.genOperations()
	if err != nil {
		return nil, err
	}
	err = e.generator.WriteLineBreak(oneLineBreak)
	if err != nil {
		return nil, err
	}

	err = e.genClientCode()
	if err != nil {
		return nil, err
	}
	err = e.generator.WriteLineBreak(oneLineBreak)
	if err != nil {
		return nil, err
	}

	return e.generator.Generate()
}

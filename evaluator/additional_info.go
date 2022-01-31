// Package evaluator contains the logic responsible for evaluating schema & query GraphQL Abstract Syntax Tree [AST].
// It orchestrates the generation of the code by using visitor & generator packages.
package evaluator

// AdditionalInfo stores all optional arguments passed to grafikgen through CLI as flags.
type AdditionalInfo struct {
	PackageName string
	ClientName  string
	UsePointers bool
}

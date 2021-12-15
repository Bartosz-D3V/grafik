package evaluator

// AdditionalInfo stores all optional arguments passed to grafikgen through CLI as flags.
type AdditionalInfo struct {
	PackageName string
	ClientName  string
	UsePointers bool
}

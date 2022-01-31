// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

// Enum represents simplified Enum in Golang AST.
// Name is the name of the enum.
// Fields is a slice of string and represents all possible values of the enum.
type Enum struct {
	Name   string
	Fields []string
}

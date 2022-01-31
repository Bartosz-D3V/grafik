// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

// Struct represents simplified Struct in Golang AST.
// Name is the name of the struct.
// Fields is a slice of TypeArg and represents struct fields.
type Struct struct {
	Name   string
	Fields []TypeField
}

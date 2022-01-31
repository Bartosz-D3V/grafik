// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

import (
	"fmt"
	"strings"
)

// TypeArg represents simplified argument in Golang AST.
// Name is the name of the argument.
// Type is type of the argument defined as string - i.e. "string", "int", "Address" etc.
type TypeArg struct {
	Name string
	Type string
}

// ExportName converts function argument name to TitleCase.
func (t TypeArg) ExportName() string {
	return strings.Title(t.Name)
}

// ExportType converts function argument type to TitleCase excluding golang primitive types.
func (t TypeArg) ExportType() TypeArg {
	const sliceTok = "[]"
	if strings.Contains(t.Type, sliceTok) {
		dim := strings.Count(t.Type, sliceTok)
		elType := strings.TrimLeft(t.Type, sliceTok)
		if isPrimitive(elType) {
			return t
		}
		return TypeArg{
			Name: t.Name,
			Type: fmt.Sprintf("%s%s", strings.Repeat(sliceTok, dim), strings.Title(elType)),
		}
	}
	if isPrimitive(t.Type) {
		return t
	}
	return TypeArg{
		Name: t.Name,
		Type: strings.Title(t.Type),
	}
}

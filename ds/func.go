// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

import (
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"strings"
)

// Func represents simplified Function in Golang AST.
// Name is the name of the Function.
// Args is a slice of TypeArg and represents function parameters.
// Type is a string and represents return type of the function - i.e. "string", "Address" etc.
// WrapperTypes is a slice of TypeArg and represents selection set in GraphQL operation.
// It is used to create wrapper struct containing all values in selection set.
type Func struct {
	Name         string
	Args         []TypeArg
	Type         string
	WrapperTypes []TypeField
}

// JoinArgsBy returns list of function arguments as concatenated string with name and type.
func (f Func) JoinArgsBy(s string) string {
	pArgs := make([]string, len(f.Args))
	for i, arg := range f.Args {
		tArg := arg.ExportType()
		pArgs[i] = fmt.Sprintf("%s %s", tArg.Name, common.SnakeCaseToCamelCase(tArg.Type))
	}

	return strings.Join(pArgs, s)
}

// ExportName converts name of the function to Title case.
func (f Func) ExportName() string {
	return strings.Title(f.Name)
}

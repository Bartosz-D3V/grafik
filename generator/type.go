package generator

import (
	"fmt"
	"strings"
)

type FuncArg struct {
	Name string
	Type string
}

type Func struct {
	Name string
	Args []FuncArg
	Type string
}

func (f Func) JoinArgsBy(s string) string {
	var pArgs []string
	for _, arg := range f.Args {
		pArgs = append(pArgs, fmt.Sprintf("%s %s", arg.Name, arg.Type))
	}

	return strings.Join(pArgs, s)
}

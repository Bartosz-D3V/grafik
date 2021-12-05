package generator

import (
	"fmt"
	"strings"
)

type TypeArg struct {
	Name string
	Type string
}

func (t TypeArg) ExportName() string {
	return strings.Title(t.Name)
}

func (t TypeArg) PointerType() string {
	if strings.Contains(t.Type, "[]") {
		return t.Type
	}
	return fmt.Sprintf("*%s", t.Type)
}

type Func struct {
	Name        string
	Args        []TypeArg
	Type        string
	WrapperArgs []TypeArg
}

func (f Func) JoinArgsBy(s string) string {
	var pArgs []string
	for _, arg := range f.Args {
		pArgs = append(pArgs, fmt.Sprintf("%s %s", arg.Name, arg.Type))
	}

	return strings.Join(pArgs, s)
}

func (f Func) ExportName() string {
	return strings.Title(f.Name)
}

type Struct struct {
	Name   string
	Fields []TypeArg
}

type Enum struct {
	Name   string
	Fields []string
}

type Const struct {
	Name string
	Val  string
}

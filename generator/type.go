package generator

import (
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"strings"
)

type TypeArg struct {
	Name string
	Type string
}

func (t TypeArg) ExportName() string {
	return strings.Title(t.Name)
}

func (t TypeArg) TitleType() TypeArg {
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

func (t TypeArg) PointerType() TypeArg {
	if strings.Contains(t.Type, "[]") {
		return t
	}
	return TypeArg{
		Name: t.Name,
		Type: fmt.Sprintf("*%s", t.Type),
	}
}

func isPrimitive(s string) bool {
	return s == "string" || s == "int" || s == "float" || s == "byte" || s == "bool"
}

type Func struct {
	Name        string
	Args        []TypeArg
	Type        string
	WrapperArgs []TypeArg
}

func (f Func) JoinArgsBy(s string) string {
	pArgs := make([]string, len(f.Args))
	for i, arg := range f.Args {
		tArg := arg.TitleType()
		pArgs[i] = fmt.Sprintf("%s %s", tArg.Name, common.SnakeCaseToCamelCase(tArg.Type))
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

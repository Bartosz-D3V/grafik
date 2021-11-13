package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

func (e *evaluator) genSchemaDef() {
	if e.schema.Query != nil {
		e.generateFromDefinition(e.schema.Query)
	}

	e.generator.WriteLineBreaks(2)

	if e.schema.Mutation != nil {
		e.generateFromDefinition(e.schema.Mutation)
	}
}

func (e *evaluator) generateFromDefinition(query *ast.Definition) {
	usrFields := query.Fields[:numOfBuiltIns(query)]
	interfaceName := query.Name

	var funcs []generator.Func
	for _, field := range usrFields {
		f := generator.Func{
			Name: field.Name,
			Args: parseArgs(field.Arguments),
			Type: field.Type.NamedType,
		}
		funcs = append(funcs, f)
	}
	e.generator.WriteInterface(interfaceName, funcs...)
}

func numOfBuiltIns(query *ast.Definition) int {
	if query.OneOf("Query") {
		return len(query.Fields) - 2
	} else {
		return len(query.Fields)
	}
}

func parseArgs(args ast.ArgumentDefinitionList) []generator.FuncArg {
	var funcArgs []generator.FuncArg
	for _, arg := range args {
		farg := generator.FuncArg{
			Name: arg.Name,
			Type: convGoType(arg.Type.NamedType),
		}
		funcArgs = append(funcArgs, farg)
	}
	return funcArgs
}

func convGoType(namedType string) string {
	switch namedType {
	case "String",
		"Int":
		return strings.ToLower(namedType)
	case "ID":
		return "string"
	case "Float":
		return "int"
	case "Boolean":
		return "bool"
	default:
		return namedType
	}
}

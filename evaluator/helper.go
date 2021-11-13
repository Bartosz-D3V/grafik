package evaluator

import (
	"github.com/Bartosz-D3V/ggrafik/generator"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

func (e *evaluator) genSchemaDef() {
	if e.schema.Query != nil {
		e.generateInterfaceFromDefinition(e.schema.Query)
	}

	e.generator.WriteLineBreak(2)

	if e.schema.Mutation != nil {
		e.generateInterfaceFromDefinition(e.schema.Mutation)
	}

	e.generator.WriteLineBreak(2)

	e.generateStructFromDefinition()
}

func (e *evaluator) generateInterfaceFromDefinition(query *ast.Definition) {
	usrFields := query.Fields[:e.numOfBuiltIns(query)]
	interfaceName := query.Name

	var funcs []generator.Func
	for _, field := range usrFields {
		f := generator.Func{
			Name: field.Name,
			Args: e.parseFnArgs(field.Arguments),
			Type: e.convGoType(field.Type.NamedType),
		}
		funcs = append(funcs, f)
	}
	e.generator.WriteInterface(interfaceName, funcs...)
}

func (e *evaluator) generateStructFromDefinition() {
	for customType := range e.customTypes {
		cType := e.schema.Types[customType]

		s := generator.Struct{
			Name:   cType.Name,
			Fields: e.parseFieldArgs(cType.Fields),
		}
		e.generator.WriteStruct(s)
	}
}

func (e *evaluator) numOfBuiltIns(query *ast.Definition) int {
	if query.OneOf("Query") {
		return len(query.Fields) - 2
	} else {
		return len(query.Fields)
	}
}

func (e *evaluator) parseFnArgs(args ast.ArgumentDefinitionList) []generator.TypeArg {
	var funcArgs []generator.TypeArg
	for _, arg := range args {
		farg := generator.TypeArg{
			Name: arg.Name,
			Type: e.convGoType(arg.Type.NamedType),
		}
		funcArgs = append(funcArgs, farg)
	}
	return funcArgs
}

func (e *evaluator) parseFieldArgs(args ast.FieldList) []generator.TypeArg {
	var funcArgs []generator.TypeArg
	for _, arg := range args {
		farg := generator.TypeArg{
			Name: arg.Name,
			Type: e.convGoType(arg.Type.NamedType),
		}
		funcArgs = append(funcArgs, farg)
	}
	return funcArgs
}

func (e *evaluator) convGoType(namedType string) string {
	switch namedType {
	case "String",
		"",
		"Int":
		return strings.ToLower(namedType)
	case "ID":
		return "string"
	case "Float":
		return "int"
	case "Boolean":
		return "bool"
	default:
		e.customTypes[namedType] = true
		return namedType
	}
}

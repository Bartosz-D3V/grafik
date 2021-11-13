package evaluator

import (
	"fmt"
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
			Type: e.convGoType(field.Type),
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
		e.generator.WriteLineBreak(2)
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
			Type: e.convGoType(arg.Type),
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
			Type: e.convGoType(arg.Type),
		}
		funcArgs = append(funcArgs, farg)
	}
	return funcArgs
}

func (e *evaluator) convGoType(astType *ast.Type) string {
	switch namedType := astType.NamedType; namedType {
	case "String",
		//"",
		"Int":
		return strings.ToLower(namedType)
	case "ID":
		return "string"
	case "Float":
		return "int"
	case "Boolean":
		return "bool"
	default:
		return e.convComplexType(astType)
	}
}

func (e *evaluator) convComplexType(astType *ast.Type) string {
	if isList := astType.IsCompatible(ast.ListType(astType.Elem, astType.Position)); isList {
		e.customTypes[astType.Elem.NamedType] = true
		return fmt.Sprintf("[]%s", astType.Elem.NamedType)
	}

	if astType.NamedType == "" {
		return ""
	}

	e.customTypes[astType.NamedType] = true
	return astType.NamedType
}

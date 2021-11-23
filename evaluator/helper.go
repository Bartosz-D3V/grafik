package evaluator

import (
	"fmt"
	"github.com/Bartosz-D3V/ggrafik/common"
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

	e.generateEnumTypesFromDefinition(e.schema.Types)

	e.generateStructs()
}

func (e *evaluator) generateInterfaceFromDefinition(query *ast.Definition) {
	usrFields := query.Fields[:common.NumOfBuiltIns(query)]
	interfaceName := query.Name

	var funcs []generator.Func
	for _, field := range usrFields {
		f := generator.Func{
			Name: field.Name,
			Args: e.parseFnArgs(&field.Arguments),
			Type: e.convGoType(field.Type),
		}
		funcs = append(funcs, f)
	}
	e.generator.WriteInterface(interfaceName, funcs...)
}

func (e *evaluator) generateEnumTypesFromDefinition(types map[string]*ast.Definition) {
	for _, definition := range types {
		if definition.Kind == ast.Enum && !definition.BuiltIn {
			e.createEnum(definition)
		}
	}
}

func (e *evaluator) generateStructs() {
	cTypes := e.visitor.IntrospectTypes()
	for _, customType := range cTypes {
		cType := e.schema.Types[customType]

		switch cType.Kind {
		case ast.Object,
			ast.InputObject:
			e.createStruct(cType)
		default:
			panic(fmt.Errorf("%s type not supported", cType.Kind))
		}
	}
}

func (e *evaluator) createEnum(cType *ast.Definition) {
	fields := make([]string, len(cType.EnumValues))
	for i, field := range cType.EnumValues {
		fields[i] = field.Name
	}

	en := generator.Enum{
		Name:   cType.Name,
		Fields: fields,
	}

	e.generator.WriteLineBreak(2)
	e.generator.WriteEnum(en)
}

func (e *evaluator) createStruct(cType *ast.Definition) {
	s := generator.Struct{
		Name:   cType.Name,
		Fields: e.parseFieldArgs(&cType.Fields),
	}
	e.generator.WriteLineBreak(2)
	e.generator.WriteStruct(s)
}

func (e *evaluator) parseFnArgs(args *ast.ArgumentDefinitionList) []generator.TypeArg {
	var funcArgs []generator.TypeArg
	for _, arg := range *args {
		farg := generator.TypeArg{
			Name: arg.Name,
			Type: e.convGoType(arg.Type),
		}
		funcArgs = append(funcArgs, farg)
	}
	return funcArgs
}

func (e *evaluator) parseFieldArgs(args *ast.FieldList) []generator.TypeArg {
	var funcArgs []generator.TypeArg
	for _, arg := range *args {
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
	if common.IsList(astType) {
		return fmt.Sprintf("[]%s", astType.Elem.NamedType)
	}

	if astType.NamedType == "" {
		return ""
	}

	return astType.NamedType
}

func (e *evaluator) genOperations() {
	ops := e.queryDocument.Operations
	opsCount := len(ops)
	src := ops[0].Position.Src.Input

	var curOp *ast.OperationDefinition
	var nextOp *ast.OperationDefinition

	for i := 0; i < opsCount; i++ {
		curOp = ops[i]
		if i < opsCount-1 {
			nextOp = ops[i+1]
		} else {
			nextOp = nil
		}
		if nextOp != nil {
			queryStr := src[curOp.Position.Start:nextOp.Position.Start]
			c := generator.Const{
				Name: curOp.Name,
				Val:  queryStr,
			}
			e.generator.WriteConst(c)
			e.generator.WriteLineBreak(2)
		} else {
			queryStr := src[curOp.Position.Start:]
			c := generator.Const{
				Name: curOp.Name,
				Val:  queryStr,
			}
			e.generator.WriteConst(c)
			e.generator.WriteLineBreak(1)
		}
	}
}

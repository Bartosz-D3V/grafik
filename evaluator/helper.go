package evaluator

import (
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

func (e *evaluator) genSchemaDef(usePointers bool) {
	e.generator.WriteLineBreak(2)

	e.generator.WriteLineBreak(2)

	e.generateEnumTypesFromDefinition(e.schema.Types)

	e.generateStructs(usePointers)

	e.generator.WriteLineBreak(2)
}

func (e *evaluator) generateEnumTypesFromDefinition(types map[string]*ast.Definition) {
	for _, definition := range types {
		if definition.Kind == ast.Enum && !definition.BuiltIn {
			e.createEnum(definition)
		}
	}
}

func (e *evaluator) generateStructs(usePointers bool) {
	cTypes := e.visitor.IntrospectTypes()
	for _, customType := range cTypes {
		cType := e.schema.Types[customType]

		switch cType.Kind {
		case ast.Object,
			ast.InputObject:
			e.createStruct(cType, usePointers)
		default:
			//panic(fmt.Errorf("%s type not supported", cType.Kind))
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

func (e *evaluator) createStruct(cType *ast.Definition, usePointers bool) {
	s := generator.Struct{
		Name:   cType.Name,
		Fields: e.parseFieldArgs(&cType.Fields),
	}
	e.generator.WriteLineBreak(2)
	e.generator.WritePublicStruct(s, usePointers)
}

func (e *evaluator) parseSelectionSet(set ast.SelectionSet) []generator.TypeArg {
	selectionSet := make([]generator.TypeArg, len(set))
	for i, s := range set {
		astField := s.(*ast.Field)
		selectionSet[i] = generator.TypeArg{
			Name: astField.Alias,
			Type: e.convGoType(astField.Definition.Type),
		}
	}
	return selectionSet
}

func (e *evaluator) parseFnArgs(args *ast.VariableDefinitionList) []generator.TypeArg {
	var funcArgs []generator.TypeArg
	for _, arg := range *args {
		farg := generator.TypeArg{
			Name: arg.Variable,
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
		if nt := astType.Elem; common.IsComplex(nt) && !common.IsList(nt) {
			return fmt.Sprintf("[]%s", nt.NamedType)
		} else {
			return fmt.Sprintf("[]%s", strings.ToLower(e.convGoType(nt)))
		}
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

func (e *evaluator) genClientCode() {
	e.genOpsInterface()
	e.generator.WriteLineBreak(2)

	e.genClientStruct()
	e.generator.WriteLineBreak(2)

	e.generator.WriteClientConstructor(e.AdditionalInfo.ClientName)
}

func (e *evaluator) genOpsInterface() {
	ops := e.queryDocument.Operations

	var funcs []generator.Func
	for _, op := range ops {
		f := generator.Func{
			Name:        op.Name,
			Args:        e.parseFnArgs(&op.VariableDefinitions),
			Type:        "(*http.Response, error)",
			WrapperArgs: e.parseSelectionSet(op.SelectionSet),
		}
		funcs = append(funcs, f)
	}
	e.generator.WriteInterface(e.AdditionalInfo.ClientName, funcs...)
	e.generator.WriteLineBreak(2)
	for _, f := range funcs {
		e.generator.WriteInterfaceImplementation(e.AdditionalInfo.ClientName, f)
		e.generator.WriteLineBreak(2)
	}

	for _, f := range funcs {
		e.genWrapperResponseStruct(f)
	}

	e.genErrorStructs()
}

func (e *evaluator) genWrapperResponseStruct(f generator.Func) {
	dataStructName := fmt.Sprintf("%sData", strings.Title(f.Name))
	responseStructName := fmt.Sprintf("%sResponse", strings.Title(f.Name))
	structWrapper := generator.Struct{
		Name: responseStructName,
		Fields: []generator.TypeArg{
			{
				Name: "data",
				Type: dataStructName,
			},
			{
				Name: "errors",
				Type: fmt.Sprintf("[]%s", generator.GraphQLErrorStructName),
			},
		},
	}
	e.generator.WritePublicStruct(structWrapper, e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(2)

	s := generator.Struct{
		Name:   dataStructName,
		Fields: f.WrapperArgs,
	}
	e.generator.WritePublicStruct(s, e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(2)
}

func (e *evaluator) genErrorStructs() {
	e.generator.WriteGraphqlErrorStructs(e.AdditionalInfo.UsePointers)
}

func (e *evaluator) genClientStruct() {
	s := generator.Struct{
		Name: e.AdditionalInfo.ClientName,
		Fields: []generator.TypeArg{
			{
				Name: "ctrl",
				Type: "graphqlClient.Client",
			},
		},
	}
	e.generator.WritePrivateStruct(s)
}

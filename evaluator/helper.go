package evaluator

import (
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

// genSchemaDef generates custom, user-defined structs and enums used in GraphQL query file.
func (e *evaluator) genSchemaDef(usePointers bool) {
	e.generator.WriteLineBreak(2)

	e.generator.WriteLineBreak(2)

	e.generateEnumTypesFromDefinition(e.schema.Types)

	e.generateStructs(usePointers)

	e.generator.WriteLineBreak(2)
}

// generateEnumTypesFromDefinition generates enums based on GraphQL schema.
func (e *evaluator) generateEnumTypesFromDefinition(types map[string]*ast.Definition) {
	for _, definition := range types {
		if definition.Kind == ast.Enum && !definition.BuiltIn {
			e.createEnum(definition)
		}
	}
}

// generateEnumTypesFromDefinition generates structs based on GraphQL schema.
func (e *evaluator) generateStructs(usePointers bool) {
	cTypes := e.visitor.IntrospectTypes()
	for _, customType := range cTypes {
		cType := e.schema.Types[customType]

		switch cType.Kind {
		case ast.Object,
			ast.InputObject:
			e.createStruct(cType, usePointers)
		}
	}
}

// createEnum creates generator.Enum and writes to IO.
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

// createStruct creates generator.Struct and writes to IO.
func (e *evaluator) createStruct(cType *ast.Definition, usePointers bool) {
	s := generator.Struct{
		Name:   cType.Name,
		Fields: e.parseFieldArgs(&cType.Fields),
	}
	e.generator.WriteLineBreak(2)
	e.generator.WritePublicStruct(s, usePointers)
}

// parseSelectionSet creates array of type generator.TypeArg based on selection set.
// Consider this GraphQL query:
// query getContinentsAndCountries {
//    continents {
//        code
//    }
//    country {
//        name
//    }
// }
// ast.SelectionSet is array of continents and country.
// parseSelectionSet will return array of type generator.TypeArg with two elements - continents and country.
// Name will be continents and country. Type will be introspected and either primitive or user defined struct.
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

// parseFnArgs converts GraphQL operation (query/mutation) arguments (ast.VariableDefinitionList) and returns slice of generator.TypeArg
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

// parseFieldArgs converts GraphQL fields (ast.FieldList) into generator.TypeArg
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

// convGoType maps GraphQL types into Go types
// User defined types are returned as they are and defined later on
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

// convComplexType recursively checks GraphQL type and returns corresponding Go type
func (e *evaluator) convComplexType(astType *ast.Type) string {
	if common.IsList(astType) {
		// If astType is not multi-dimensional array return '[]' with named type
		if nt := astType.Elem; common.IsComplex(nt) && !common.IsList(nt) {
			return fmt.Sprintf("[]%s", nt.NamedType)
		} else {
			return fmt.Sprintf("[]%s", e.convGoType(nt))
		}
	}

	// If astType is not array it is an object - just return the name
	return astType.NamedType
}

// genOperations generates GraphQL operations as constants.
// For example the following query:
// query getContinentsAndCountries {
//    continents {
//        code
//    }
//    country {
//        name
//    }
// }
// Will be conversed to this Go code:
// const getContinentsAndCountries = `query getContinentsAndCountries {
//    continents {
//        code
//    }
//    country {
//        name
//    }
// }'
func (e *evaluator) genOperations() {
	ops := e.queryDocument.Operations
	opsCount := len(ops)
	src := ops[0].Position.Src.Input

	var curOp *ast.OperationDefinition
	var nextOp *ast.OperationDefinition

	// Split multiple GraphQL operations into sub-operations to generate const value for each operation
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

// genClientCode generates client code - all interfaces, constructor methods and client struct.
func (e *evaluator) genClientCode() {
	e.genOpsInterface()
	e.generator.WriteLineBreak(2)

	e.genClientStruct()
	e.generator.WriteLineBreak(2)

	e.generator.WriteClientConstructor(e.AdditionalInfo.ClientName)
}

// genOpsInterface generates public interface for grafik client.
// For example:
// type SpaceXClient interface {
//	AddOrUpdateHardcodedUser(rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error)
// }
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

	// Generate interface implementation for each interface method
	for _, f := range funcs {
		e.generator.WriteInterfaceImplementation(e.AdditionalInfo.ClientName, f)
		e.generator.WriteLineBreak(2)
	}

	// Generate wrapper struct for selection set operations
	for _, f := range funcs {
		e.genWrapperResponseStruct(f)
	}

	// Generate predefined error structs
	e.genErrorStructs()
}

// genWrapperResponseStruct generates top level GraphQL response type
// See https://graphql.org/learn/serving-over-http/#response
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

	// generate object referenced in 'data' JSON response
	// if object has selection set - those will be created as struct fields
	s := generator.Struct{
		Name:   dataStructName,
		Fields: f.WrapperArgs,
	}
	e.generator.WritePublicStruct(s, e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(2)
}

// genErrorStructs generates predefined GraphQL error structs
func (e *evaluator) genErrorStructs() {
	e.generator.WriteGraphqlErrorStructs(e.AdditionalInfo.UsePointers)
}

// genClientStruct generates internal grafik GraphQL client defined in package client
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

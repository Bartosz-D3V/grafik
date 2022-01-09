package evaluator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/Bartosz-D3V/grafik/generator"
	"github.com/vektah/gqlparser/ast"
	"sort"
	"strings"
)

const (
	oneLineBreak              = 1
	twoLinesBreak             = 2
	graphQLFragmentStructName = "Fragment"
	graphQLUnionStructName    = "Union"
)

// genSchemaDef generates custom, user-defined structs and enums used in GraphQL query file.
func (e *evaluator) genSchemaDef() {
	e.generator.WriteLineBreak(twoLinesBreak)

	e.generateGoTypes()

	e.generator.WriteLineBreak(twoLinesBreak)
}

// generateGoTypes iterates through all fields in GraphQL query and generates GO type based on selected subfields.
func (e *evaluator) generateGoTypes() {
	cTypes := e.visitor.IntrospectTypes()

	// To make the output order of the generated code deterministic always sort alphabetically.
	keys := make([]string, 0, len(cTypes))
	for k := range cTypes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		cType, ok := e.schema.Types[key]
		if !ok {
			panic(fmt.Errorf("failed to find definition of %s in GraphQL Schema AST", key))
		}

		switch cType.Kind {
		case ast.Object,
			ast.InputObject:
			e.createStruct(cType, cTypes[key])
		case ast.Enum:
			e.createEnum(cType)
		case ast.Scalar:
			e.createInterfaceType(cType)
		case ast.Interface:
			e.createCommonStruct(cType, cTypes[key], graphQLFragmentStructName)
		case ast.Union:
			e.createCommonStruct(cType, cTypes[key], graphQLUnionStructName)
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

	e.generator.WriteLineBreak(twoLinesBreak)
	e.generator.WriteEnum(en)
}

// createInterface creates type 'any' in Go [type X interface{}] and writes to IO.
func (e *evaluator) createInterfaceType(cType *ast.Definition) {
	e.generator.WriteLineBreak(twoLinesBreak)
	e.generator.WriteInterface(cType.Name)
}

// createStruct creates generator.Struct and writes to IO.
func (e *evaluator) createStruct(cType *ast.Definition, selectedFields []string) {
	s := generator.Struct{
		Name:   cType.Name,
		Fields: e.parseFieldArgs(&cType.Fields, selectedFields),
	}
	e.generator.WriteLineBreak(twoLinesBreak)
	e.generator.WritePublicStruct(s, e.AdditionalInfo.UsePointers)
}

// createCommonStruct creates a generic struct containing all the fields that interface and all implementations it has.
func (e *evaluator) createCommonStruct(cType *ast.Definition, selectedFields []string, graphQLTypeSuffix string) {
	fragmentName := fmt.Sprintf("%s%s", cType.Name, graphQLTypeSuffix)
	fragmentFields := make(ast.FieldList, 0)

	// Add fields of all implementations
	for _, definition := range e.schema.GetPossibleTypes(cType) {
		fragmentFields = append(fragmentFields, definition.Fields...)
	}

	// Add fields of an interface
	fragmentFields = append(fragmentFields, cType.Fields...)

	fList := make(ast.FieldList, 0)
	for _, fField := range fragmentFields {
		selected := false
		for _, field := range selectedFields {
			if fField.Name == field {
				selected = true
			}
		}
		if !selected {
			continue
		}

		if fList.ForName(fField.Name) == nil {
			fList = append(fList, fField)
		}
	}

	fragmentDef := &ast.Definition{
		Kind:   ast.Object,
		Name:   fragmentName,
		Fields: fList,
	}

	allFields := make([]string, len(fList))
	for i, field := range fList {
		allFields[i] = field.Name
	}

	e.createStruct(fragmentDef, allFields)
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

// parseFnArgs converts GraphQL operation (query/mutation) arguments (ast.VariableDefinitionList) and returns slice of generator.TypeArg.
func (e *evaluator) parseFnArgs(args *ast.VariableDefinitionList) []generator.TypeArg {
	funcArgs := make([]generator.TypeArg, len(*args))
	for i, arg := range *args {
		fArg := generator.TypeArg{
			Name: arg.Variable,
			Type: e.convGoType(arg.Type),
		}
		funcArgs[i] = fArg
	}
	return funcArgs
}

// parseFieldArgs converts GraphQL fields (ast.FieldList) into generator.TypeArg.
func (e *evaluator) parseFieldArgs(args *ast.FieldList, selectedFields []string) []generator.TypeArg {
	funcArgs := make([]generator.TypeArg, 0)
	for _, arg := range *args {
		selected := false
		for _, field := range selectedFields {
			if arg.Name == field {
				selected = true
			}
		}
		if !selected {
			continue
		}

		fArg := generator.TypeArg{
			Name: arg.Name,
			Type: e.convGoType(arg.Type),
		}
		funcArgs = append(funcArgs, fArg)
	}
	return funcArgs
}

// convGoType maps GraphQL types into Go types.
// User defined types are returned as they are and defined later on.
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

// convComplexType recursively checks GraphQL type and returns corresponding Go type.
func (e *evaluator) convComplexType(astType *ast.Type) string {
	switch {
	case common.IsList(astType):
		return e.convListType(astType)
	// If astType is interface return with 'Fragment' suffix
	case e.schema.Types[astType.NamedType].Kind == ast.Interface:
		return fmt.Sprintf("%s%s", astType.NamedType, graphQLFragmentStructName)
	// If astType is union return with 'Union' suffix
	case e.schema.Types[astType.NamedType].Kind == ast.Union:
		return fmt.Sprintf("%s%s", astType.NamedType, graphQLUnionStructName)
	// Otherwise, just return the name
	default:
		return astType.NamedType
	}
}

// convListType returns Go type for list of GraphQL Type.
// I.e. [[Character]] -> [][]Character
func (e *evaluator) convListType(astType *ast.Type) string {
	switch nt := astType.Elem; {
	// If astType is not multi-dimensional array of interfaces return '[]' with 'Fragment' suffix
	case common.IsComplex(nt) && !common.IsList(nt) && e.schema.Types[nt.NamedType].Kind == ast.Interface:
		return fmt.Sprintf("[]%s%s", nt.NamedType, graphQLFragmentStructName)
	// If astType is not multi-dimensional array of unions return '[]' with 'Union' suffix
	case common.IsComplex(nt) && !common.IsList(nt) && e.schema.Types[nt.NamedType].Kind == ast.Union:
		return fmt.Sprintf("[]%s%s", nt.NamedType, graphQLUnionStructName)
	// If astType is not multi-dimensional array return '[]' with named type
	case common.IsComplex(nt) && !common.IsList(nt):
		return fmt.Sprintf("[]%s", nt.NamedType)
	// Otherwise, recursively check the type
	default:
		return fmt.Sprintf("[]%s", e.convGoType(nt))
	}
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
			queryStr = e.removeComments(queryStr)
			c := generator.Const{
				Name: curOp.Name,
				Val:  queryStr,
			}
			e.generator.WriteConst(c)
			e.generator.WriteLineBreak(twoLinesBreak)
		} else {
			queryStr := src[curOp.Position.Start:]
			queryStr = e.removeComments(queryStr)
			c := generator.Const{
				Name: curOp.Name,
				Val:  queryStr,
			}
			e.generator.WriteConst(c)
			e.generator.WriteLineBreak(oneLineBreak)
		}
	}
}

// genClientCode generates client code - all interfaces, constructor methods and client struct.
func (e *evaluator) genClientCode() {
	e.genOpsInterface()
	e.generator.WriteLineBreak(twoLinesBreak)

	e.genClientStruct()
	e.generator.WriteLineBreak(twoLinesBreak)

	e.generator.WriteClientConstructor(e.AdditionalInfo.ClientName)
}

// genOpsInterface generates public interface for grafik client.
// For example:
// type SpaceXClient interface {
//	AddOrUpdateHardcodedUser(rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error)
// }
func (e *evaluator) genOpsInterface() {
	ops := e.queryDocument.Operations

	funcs := make([]generator.Func, len(ops))
	for i, op := range ops {
		f := generator.Func{
			Name:         op.Name,
			Args:         e.parseFnArgs(&op.VariableDefinitions),
			Type:         "(*http.Response, error)",
			WrapperTypes: e.parseSelectionSet(op.SelectionSet),
		}
		funcs[i] = f
	}
	e.generator.WriteInterface(e.AdditionalInfo.ClientName, funcs...)
	e.generator.WriteLineBreak(twoLinesBreak)

	// Generate interface implementation for each interface method
	for _, f := range funcs {
		e.generator.WriteInterfaceImplementation(e.AdditionalInfo.ClientName, f)
		e.generator.WriteLineBreak(twoLinesBreak)
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
	e.generator.WriteLineBreak(twoLinesBreak)

	// generate object referenced in 'data' JSON response
	// if object has selection set - those will be created as struct fields
	s := generator.Struct{
		Name:   dataStructName,
		Fields: f.WrapperTypes,
	}
	e.generator.WritePublicStruct(s, e.AdditionalInfo.UsePointers)
	e.generator.WriteLineBreak(twoLinesBreak)
}

// genErrorStructs generates predefined GraphQL error structs.
func (e *evaluator) genErrorStructs() {
	e.generator.WriteGraphqlErrorStructs(e.AdditionalInfo.UsePointers)
}

// genClientStruct generates internal grafik GraphQL client defined in package client.
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

// removeComments removes comments from GraphQL queries.
func (e *evaluator) removeComments(queryStr string) string {
	const commentToken = "#"

	var out bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(queryStr))

	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), commentToken)
		code := strings.TrimSpace(splitLine[0])
		if code != "" {
			out.WriteString(splitLine[0])
			out.WriteRune('\n')
		}
	}
	return strings.TrimRight(out.String(), "\n")
}

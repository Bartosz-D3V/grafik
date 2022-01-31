// Package visitor abstracts logic responsible for determining which custom types from GraphQL Schema file should be generated based on usage in GraphQL query file.
package visitor

import (
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/vektah/gqlparser/ast"
)

// parseOpTypes parses selectionSet of each GraphQL operation and all variables.
func (v *visitor) parseOpTypes(opList ast.OperationList) {
	for _, opDef := range opList {
		v.parseSelectionSet(opDef.SelectionSet, make([]string, 0), false)
		v.parseVariables(opDef.VariableDefinitions)
	}
}

// parseSelectionSet parses each selection based on its type (Field/FragmentSpread/Inline Fragment)
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseSelectionSet(selectionSet ast.SelectionSet, fields []string, registerType bool) []string {
	for _, selection := range selectionSet {
		switch selectionType := selection.(type) {
		case *ast.Field:
			fields = append(fields, selectionType.Name)
			v.parseSelectionSet(selectionType.SelectionSet, make([]string, 0), true)
			if registerType {
				v.registerTypeByName(selectionType.ObjectDefinition.Name, fields)
				v.registerType(selectionType.Definition.Type, make([]string, 0))
			}
		case *ast.InlineFragment:
			fields = v.parseInlineFragment(selectionType, make([]string, 0), false)
			if registerType {
				v.registerTypeByName(selectionType.ObjectDefinition.Name, fields)
			}
		case *ast.FragmentSpread:
			fields = v.parseFragmentSpread(selectionType, fields, true)
		}
	}
	return fields
}

// parseFragmentSpread parses GraphQL Fragment Spread - it will add all fields to the visitor.
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseFragmentSpread(fragmentSpread *ast.FragmentSpread, fields []string, registerType bool) []string {
	fields = v.parseSelectionSet(fragmentSpread.Definition.SelectionSet, fields, registerType)
	return fields
}

// parseInlineFragment parses GraphQL Inline Fragment - it will add all fields of all fragments to the visitor.
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseInlineFragment(parsedType *ast.InlineFragment, fields []string, registerType bool) []string {
	fields = v.parseSelectionSet(parsedType.SelectionSet, fields, registerType)
	return fields
}

// parseVariables parses all variables defined in the GraphQL operation.
func (v *visitor) parseVariables(variableDefinitionList ast.VariableDefinitionList) {
	for _, varDef := range variableDefinitionList {
		v.parseType(varDef.Type)
	}
}

// parseType parses generic GraphQL Type.
func (v *visitor) parseType(astType *ast.Type) {
	leafType := v.findLeafType(astType)
	leafTypeDef := v.schema.Types[leafType.NamedType]

	// If the type is not built-in to the GraphQL specification, register it with all fields selected in the GraphQL query.
	if leafTypeDef != nil && !leafTypeDef.BuiltIn {
		fields := make([]string, len(leafTypeDef.Fields))
		for i, field := range leafTypeDef.Fields {
			fields[i] = field.Name
			v.parseType(field.Type)
		}
		v.registerType(astType, fields)
	}
}

// registerType adds field with selected fields into visitor.
func (v *visitor) registerType(astType *ast.Type, fields []string) {
	leafType := v.findLeafType(astType)

	if leafType == nil || v.schema.Types[leafType.NamedType].BuiltIn {
		return
	}

	if cFields, ok := v.customTypes[leafType.NamedType]; ok {
		fields = append(cFields, fields...)
		v.customTypes[leafType.NamedType] = fields
	} else {
		v.customTypes[leafType.NamedType] = fields
	}
}

// registerTypeByName adds field by name with selected fields into visitor.
func (v *visitor) registerTypeByName(astTypeName string, fields []string) {
	if v.schema.Types[astTypeName].BuiltIn {
		return
	}

	if cFields, ok := v.customTypes[astTypeName]; ok {
		fields = append(cFields, fields...)
		v.customTypes[astTypeName] = fields
	} else {
		v.customTypes[astTypeName] = fields
	}
}

// findLeafType unwraps the type of array.
// If the type is a list (i.e. [[Character!]]) then return leafType (in this example Character).
func (v *visitor) findLeafType(astType *ast.Type) *ast.Type {
	if common.IsList(astType) {
		return v.findLeafType(astType.Elem)
	}
	return astType
}

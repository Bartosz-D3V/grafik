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
//func (v *visitor) parseSelectionSet(selectionSet ast.SelectionSet, fields []string) []string {
//	for _, selection := range selectionSet {
//		switch selectionType := selection.(type) {
//		case *ast.Field:
//			//fields = v.parseField(selectionType, fields)
//			fields = append(fields, selectionType.Name)
//			v.parseSelectionSet(selectionType.SelectionSet,make([]string,0))
//			v.registerType2(selectionType.ObjectDefinition.Name,fields)
//		case *ast.InlineFragment:
//			fields = v.parseInlineFragment(selectionType, make([]string,0))
//			v.registerType2(selectionType.ObjectDefinition.Name, fields)
//		case *ast.FragmentSpread:
//			fields = v.parseFragmentSpread(selectionType, fields)
//		}
//	}
//	return fields
//}

// TODO Trial
func (v *visitor) parseSelectionSet(selectionSet ast.SelectionSet, fields []string, registerType bool) []string {
	for _, selection := range selectionSet {
		switch selectionType := selection.(type) {
		case *ast.Field:
			//fields = v.parseField(selectionType, fields)
			fields = append(fields, selectionType.Name)
			v.parseSelectionSet(selectionType.SelectionSet, make([]string, 0), true)
			if registerType {
				v.registerType2(selectionType.ObjectDefinition.Name, fields)
				v.registerType(selectionType.Definition.Type, make([]string, 0))
			}
		case *ast.InlineFragment:
			fields = v.parseInlineFragment(selectionType, make([]string, 0), false)
			if registerType {
				v.registerType2(selectionType.ObjectDefinition.Name, fields)
			}
		case *ast.FragmentSpread:
			fields = v.parseFragmentSpread(selectionType, fields, true)
		}
	}
	return fields
}

// parseField parses GraphQL field and registers its type.
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseField(field *ast.Field, fields []string) []string {
	//if field.SelectionSet == nil || len(field.SelectionSet) == 0 {
	//	fields = append(fields, field.Name)
	//}
	fields = append(fields, field.Name)
	//for _, s := range field.SelectionSet {
	//	fields = v.parseSelection(s, fields)
	//}

	//if field.ObjectDefinition.Interfaces != nil {
	//	for _, i := range field.ObjectDefinition.Interfaces {
	//		v.registerType2(i, fields)
	//	}
	//} else {
	//	v.registerType(field.Definition.Type, fields)
	//}

	// If the fields is a selectionSet - parse it recursively.
	//if field.SelectionSet != nil && len(field.SelectionSet) > 0 {
	//	v.parseSelectionSet(field.SelectionSet, fields)
	//}
	return fields
}

// parseSelection parses GraphQL selection.
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseSelection(s ast.Selection, fields []string) []string {
	switch parsedType := s.(type) {
	case *ast.Field:
		//fields = v.parseField(parsedType, fields)
		fields = append(fields, parsedType.Name)
		//v.parseSelectionSet(parsedType.SelectionSet, fields)
		v.registerType(parsedType.Definition.Type, fields)
	case *ast.InlineFragment:
		fields = v.parseInlineFragment(parsedType, fields, false)
		//v.parseSelectionSet(parsedType.SelectionSet, fields)
		v.registerType2(parsedType.ObjectDefinition.Name, fields)
	case *ast.FragmentSpread:
		fields = v.parseFragmentSpread(parsedType, fields, false)
		//v.parseSelectionSet(parsedType.Definition.SelectionSet, fields)
	}
	return fields
}

// parseFragmentSpread parses GraphQL Fragment Spread - it will add all fields to the visitor.
// It returns the fields that the selection uses from GraphQL schema.
func (v *visitor) parseFragmentSpread(parsedType *ast.FragmentSpread, fields []string, registerType bool) []string {
	fields = v.parseSelectionSet(parsedType.Definition.SelectionSet, fields, registerType)
	return fields
}

// parseFragmentSpread parses GraphQL Fragment Spread - it will add all fields of all fragments to the visitor.
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
func (v *visitor) parseType(definitionType *ast.Type) {
	leafType := v.findLeafType(definitionType)
	leafTypeDef := v.schema.Types[leafType.NamedType]

	// If the type is not built-in to the GraphQL specification, register it will all fields selected in the GraphQL query.
	if leafTypeDef != nil && !leafTypeDef.BuiltIn {
		fields := make([]string, len(leafTypeDef.Fields))
		for i, field := range leafTypeDef.Fields {
			fields[i] = field.Name
			v.parseType(field.Type)
		}
		v.registerType(definitionType, fields)
	}
}

// registerType adds field with selected fields into visitor.
func (v *visitor) registerType(selectType *ast.Type, fields []string) {
	leafType := v.findLeafType(selectType)

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

func (v *visitor) registerType2(selectType string, fields []string) {
	if v.schema.Types[selectType].BuiltIn {
		return
	}

	if cFields, ok := v.customTypes[selectType]; ok {
		fields = append(cFields, fields...)
		v.customTypes[selectType] = fields
	} else {
		v.customTypes[selectType] = fields
	}
}

// findLeafType unwraps the type of array.
// If the type is a list (i.e. [[Character!]]) then return leafType (in this example Character).
func (v *visitor) findLeafType(elem *ast.Type) *ast.Type {
	if common.IsList(elem) {
		return v.findLeafType(elem.Elem)
	}
	return elem
}

package visitor

import (
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/vektah/gqlparser/ast"
)

func (v *visitor) parseOpTypes(opList ast.OperationList) {
	for _, opDef := range opList {
		v.parseSelectionSet(opDef.SelectionSet, make([]string, 0))
		v.parseVariables(opDef.VariableDefinitions)
	}
}

func (v *visitor) parseSelectionSet(selectionSet ast.SelectionSet, fields []string) []string {
	for _, selection := range selectionSet {
		switch selectionType := selection.(type) {
		case *ast.Field:
			fields = v.parseField(selectionType, fields)
		case *ast.FragmentSpread:
			fields = v.parseFragmentSpread(selectionType, fields)
		case *ast.InlineFragment:
			fields = v.parseInlineFragment(selectionType, fields)
		}
	}
	return fields
}

func (v *visitor) parseField(field *ast.Field, fields []string) []string {
	if field.SelectionSet == nil || len(field.SelectionSet) == 0 {
		fields = append(fields, field.Name)
	}

	for _, s := range field.SelectionSet {
		fields = v.parseSelection(s, fields)
	}

	v.registerType(field.Definition.Type, fields)
	if field.SelectionSet != nil && len(field.SelectionSet) > 0 {
		v.parseSelectionSet(field.SelectionSet, fields)
	}
	return fields
}

func (v *visitor) parseSelection(s ast.Selection, fields []string) []string {
	switch parsedType := s.(type) {
	case *ast.Field:
		fields = append(fields, parsedType.Name)
		v.parseSelectionSet(parsedType.SelectionSet, fields)
		v.registerType(parsedType.Definition.Type, make([]string, 0))
	case *ast.InlineFragment:
		fields = v.parseInlineFragment(parsedType, fields)
		v.parseSelectionSet(parsedType.SelectionSet, fields)
	case *ast.FragmentSpread:
		fields = v.parseFragmentSpread(parsedType, fields)
		v.parseSelectionSet(parsedType.Definition.SelectionSet, fields)
	}
	return fields
}

func (v *visitor) parseFragmentSpread(parsedType *ast.FragmentSpread, fields []string) []string {
	for _, sel := range parsedType.Definition.SelectionSet {
		switch selType := sel.(type) {
		case *ast.Field:
			fields = append(fields, selType.Name)
		case *ast.FragmentSpread:
			fields = v.parseSelectionSet(selType.Definition.SelectionSet, fields)
		case *ast.InlineFragment:
			fields = v.parseSelectionSet(selType.SelectionSet, fields)
		}
	}
	return fields
}

func (v *visitor) parseInlineFragment(parsedType *ast.InlineFragment, fields []string) []string {
	for _, sel := range parsedType.SelectionSet {
		switch selType := sel.(type) {
		case *ast.Field:
			fields = append(fields, selType.Name)
		case *ast.FragmentSpread:
			fields = v.parseSelectionSet(selType.Definition.SelectionSet, fields)
		case *ast.InlineFragment:
			fields = v.parseSelectionSet(selType.SelectionSet, fields)
		}
	}
	return fields
}

func (v *visitor) parseVariables(variableDefinitionList ast.VariableDefinitionList) {
	for _, varDef := range variableDefinitionList {
		v.parseType(varDef.Type)
	}
}

func (v *visitor) parseType(definitionType *ast.Type) {
	var leafDefType *ast.Definition

	if common.IsList(definitionType) {
		leafType := v.findLeafType(definitionType)
		leafDefType = v.schema.Types[leafType.NamedType]
	} else {
		leafDefType = v.schema.Types[definitionType.NamedType]
	}

	if leafDefType != nil && !leafDefType.BuiltIn {
		fields := make([]string, len(leafDefType.Fields))
		for i, field := range leafDefType.Fields {
			fields[i] = field.Name
			v.parseType(field.Type)
		}
		v.registerType(definitionType, fields)
	}
}

func (v *visitor) registerType(field *ast.Type, fields []string) {
	var leafType *ast.Type

	if common.IsList(field) {
		leafType = v.findLeafType(field)
	} else {
		leafType = field
	}

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

func (v *visitor) findLeafType(elem *ast.Type) *ast.Type {
	if common.IsList(elem) {
		return v.findLeafType(elem.Elem)
	}
	return elem
}

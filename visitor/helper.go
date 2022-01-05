package visitor

import (
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/vektah/gqlparser/ast"
)

func (v *visitor) parseOpTypes(opList ast.OperationList) {
	for _, opDef := range opList {
		v.parseSelection(opDef.SelectionSet, make([]string, 0))
		v.parseVariables(opDef.VariableDefinitions)
	}
}

func (v *visitor) parseSelection(selectionSet ast.SelectionSet, fields []string) []string {
	for _, selection := range selectionSet {

		if field, ok := selection.(*ast.Field); ok {
			if field.SelectionSet == nil || len(field.SelectionSet) == 0 {
				fields = append(fields, field.Name)
				continue
			}

			for _, s := range field.SelectionSet {
				switch parsedType := s.(type) {
				case *ast.Field:
					fields = append(fields, parsedType.Name)
					if parsedType.SelectionSet != nil && len(parsedType.SelectionSet) > 0 {
						v.parseSelection(parsedType.SelectionSet, fields)
					}
					v.registerType(parsedType.Definition.Type, make([]string, 0))
				case *ast.InlineFragment:
					v.parseSelection(parsedType.SelectionSet, fields)
				case *ast.FragmentSpread:
					fields = v.parseFragmentSpread(parsedType, fields)
					v.parseSelection(parsedType.Definition.SelectionSet, fields)
				}
			}

			v.registerType(field.Definition.Type, fields)
			if field.SelectionSet != nil && len(field.SelectionSet) > 0 {
				v.parseSelection(field.SelectionSet, fields)
			}
		}
		if fragment, ok := selection.(*ast.FragmentSpread); ok {
			fields = v.parseFragmentSpread(fragment, fields)
			v.parseSelection(fragment.Definition.SelectionSet, fields)
		}
	}
	return fields
}

func (v *visitor) parseFragmentSpread(parsedType *ast.FragmentSpread, fields []string) []string {
	for _, sel := range parsedType.Definition.SelectionSet {
		switch selType := sel.(type) {
		case *ast.Field:
			fields = append(fields, selType.Name)
		case *ast.FragmentSpread:
			fields = v.parseSelection(selType.Definition.SelectionSet, fields)
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
	if common.IsList(field) {
		leafType := v.findLeafType(field)
		if v.schema.Types[leafType.NamedType].BuiltIn {
			return
		}

		cFields, ok := v.customTypes[leafType.NamedType]
		if ok {
			fields = append(cFields, fields...)
			v.customTypes[leafType.NamedType] = fields
		} else {
			v.customTypes[leafType.NamedType] = fields
		}
	} else {
		if v.schema.Types[field.NamedType].BuiltIn {
			return
		}

		cFields, ok := v.customTypes[field.NamedType]
		if ok {
			fields = append(cFields, fields...)
			v.customTypes[field.NamedType] = fields
		} else {
			v.customTypes[field.NamedType] = fields
		}
	}
}

func (v *visitor) findLeafType(elem *ast.Type) *ast.Type {
	if common.IsList(elem) {
		return v.findLeafType(elem.Elem)
	}
	return elem
}

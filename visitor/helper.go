package visitor

import (
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/vektah/gqlparser/ast"
)

//func (v *visitor) parseOpTypes(query *ast.Definition) {
//	for _, field := range query.Fields {
//		if field.Type.NamedType != "" {
//			v.findSubTypes(v.schema.Types[field.Type.NamedType])
//		}
//
//		for _, arg := range field.Arguments {
//			v.findSubTypes(v.schema.Types[arg.Type.NamedType])
//		}
//
//		if common.IsList(field.Type) {
//			typeLeafType := v.findLeafType(field.Type.Elem)
//			v.findSubTypes(v.schema.Types[typeLeafType.NamedType])
//
//			for _, arg := range field.Arguments {
//				argLeafType := v.findLeafType(arg.Type)
//				v.findSubTypes(v.schema.Types[argLeafType.NamedType])
//			}
//		}
//	}
//}
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
			//fields := make([]string, 0)
			for _, s := range field.SelectionSet {
				switch parsedType := s.(type) {
				case *ast.Field:
					fields = append(fields, parsedType.Name)
					if parsedType.SelectionSet != nil && len(parsedType.SelectionSet) > 0 {
						v.parseSelection(parsedType.SelectionSet, fields)
					}
					v.registerType2(parsedType.Definition.Type, make([]string, 0))
				case *ast.InlineFragment:
					v.parseSelection(parsedType.SelectionSet, fields)
				case *ast.FragmentSpread:
					fields = v.parseFragmentSpread(parsedType, fields)
					v.parseSelection(parsedType.Definition.SelectionSet, fields)
				}
			}

			v.registerType2(field.Definition.Type, fields)
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
		v.registerType2(definitionType, fields)
	}
}

func (v *visitor) registerType2(field *ast.Type, fields []string) {
	if common.IsList(field) {
		leafType := v.findLeafType(field)
		cFields, ok := v.customTypes2[leafType.NamedType]
		if ok {
			joinedFields := append(cFields, fields...)
			v.customTypes2[leafType.NamedType] = joinedFields
		} else {
			v.customTypes2[leafType.NamedType] = fields
		}
	} else {
		cFields, ok := v.customTypes2[field.NamedType]
		if ok {
			joinedFields := append(cFields, fields...)
			v.customTypes2[field.NamedType] = joinedFields
		} else {
			v.customTypes2[field.NamedType] = fields
		}
	}
}
func (v *visitor) findSubTypes(t *ast.Definition) {
	if t != nil && !t.BuiltIn {
		v.registerType(t.Name)
		for _, field := range t.Fields {
			v.registerType(t.Name)
			v.findSubTypes(v.schema.Types[field.Type.NamedType])
			if common.IsList(field.Type) {
				v.registerType(field.Type.Elem.NamedType)
			}
		}
	}
}

func (v *visitor) findLeafType(elem *ast.Type) *ast.Type {
	if common.IsList(elem) {
		return v.findLeafType(elem.Elem)
	}
	return elem
}

func (v *visitor) registerType(typeName string) {
	if typeName != "" && !v.typeRegistered(typeName) {
		v.customTypes = append(v.customTypes, typeName)
	}
}

func (v *visitor) typeRegistered(typeName string) bool {
	for _, cType := range v.customTypes {
		if typeName == cType {
			return true
		}
	}
	return false
}

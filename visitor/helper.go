package visitor

import (
	"github.com/Bartosz-D3V/grafik/common"
	"github.com/vektah/gqlparser/ast"
)

func (v *visitor) IntrospectTypes() []string {
	if v.schema.Query != nil {
		v.parseOpTypes(v.schema.Query)
	}
	if v.schema.Mutation != nil {
		v.parseOpTypes(v.schema.Mutation)
	}
	return v.customTypes
}

func (v *visitor) parseOpTypes(query *ast.Definition) {
	for _, field := range query.Fields {
		if field.Type.NamedType != "" {
			v.findSubTypes(v.schema.Types[field.Type.NamedType])
		}

		for _, arg := range field.Arguments {
			v.findSubTypes(v.schema.Types[arg.Type.NamedType])
		}

		if common.IsList(field.Type) {
			typeLeafType := v.findLeafType(field.Type.Elem)
			v.findSubTypes(v.schema.Types[typeLeafType.NamedType])

			for _, arg := range field.Arguments {
				argLeafType := v.findLeafType(arg.Type)
				v.findSubTypes(v.schema.Types[argLeafType.NamedType])
			}
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

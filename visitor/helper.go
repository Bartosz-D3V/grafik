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
	usrFields := common.FilterCustomFields(query.Fields)
	for _, field := range usrFields {
		v.findSubTypes(v.schema.Types[field.Type.NamedType])
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
	if t != nil && t.Fields != nil {
		for _, field := range t.Fields {
			if common.IsComplex(field.Type) {
				v.registerType(t.Name)
				if !v.typeRegistered(field.Type.NamedType) {
					v.findSubTypes(v.schema.Types[field.Type.NamedType])
				}
				if common.IsList(field.Type) && common.IsComplex(field.Type.Elem) {
					v.registerType(field.Type.Elem.NamedType)
				}
			} else {
				v.registerType(t.Name)
			}
		}
	}
}

func (v *visitor) findLeafType(elem *ast.Type) *ast.Type {
	if common.IsList(elem) {
		return v.findLeafType(elem.Elem)
	} else {
		return elem
	}
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

package visitor

import (
	"github.com/Bartosz-D3V/ggrafik/common"
	"github.com/vektah/gqlparser/ast"
)

func (v *visitor) IntrospectTypes() []string {
	if v.schema.Query != nil {
		v.parseOpTypes(v.schema.Query)
	}
	return v.customTypes
}

func (v *visitor) parseOpTypes(query *ast.Definition) {
	usrFields := query.Fields[:common.NumOfBuiltIns(query)]
	for _, field := range usrFields {
		v.findSubTypes(v.schema.Types[field.Type.NamedType])
		for _, arg := range field.Arguments {
			v.findSubTypes(v.schema.Types[arg.Type.NamedType])
		}
	}
}

func (v *visitor) findSubTypes(t *ast.Definition) {
	if t != nil && t.Fields != nil {
		for _, field := range t.Fields {
			if v.isComplex(field.Type) {
				v.registerType(t.Name)
				if !v.typeRegistered(field.Type.NamedType) {
					v.findSubTypes(v.schema.Types[field.Type.NamedType])
				}
				if common.IsList(field.Type) {
					v.registerType(field.Type.Elem.NamedType)
				}
			} else {
				v.registerType(t.Name)
			}
		}
	}
}

func (v *visitor) isComplex(t *ast.Type) bool {
	return t.NamedType != "String" && t.NamedType != "Int" &&
		t.NamedType != "ID" && t.NamedType != "Float" &&
		t.NamedType != "Boolean"
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

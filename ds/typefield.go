// Package ds (Data Structure) contains all golang data structures used by generator.
package ds

import (
	"fmt"
	"strings"
)

// TypeField represents simplified struct field in Golang AST.
// Name is the name of the field.
// Type is type of the field defined as string - i.e. "string", "int", "Address" etc.
// JsonName is the name of the field used in `json:` tag.
type TypeField struct {
	Name     string
	Type     string
	JsonName string
}

// ExportName converts field name to TitleCase.
func (t TypeField) ExportName() string {
	return strings.Title(t.Name)
}

// ExportType converts field name to TitleCase excluding golang primitive types.
func (t TypeField) ExportType() TypeField {
	const sliceTok = "[]"
	if strings.Contains(t.Type, sliceTok) {
		dim := strings.Count(t.Type, sliceTok)
		elType := strings.TrimLeft(t.Type, sliceTok)
		if isPrimitive(elType) {
			return t
		}
		return TypeField{
			Name:     t.Name,
			Type:     fmt.Sprintf("%s%s", strings.Repeat(sliceTok, dim), strings.Title(elType)),
			JsonName: t.JsonName,
		}
	}
	if isPrimitive(t.Type) {
		return t
	}
	return TypeField{
		Name:     t.Name,
		Type:     strings.Title(t.Type),
		JsonName: t.JsonName,
	}
}

// PointerType converts TypeField to pointer type, excluding arrays/slices/maps.
func (t TypeField) PointerType() TypeField {
	if strings.Contains(t.Type, "[]") {
		return t
	}
	return TypeField{
		Name:     t.Name,
		Type:     fmt.Sprintf("*%s", t.Type),
		JsonName: t.JsonName,
	}
}

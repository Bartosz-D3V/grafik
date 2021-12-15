// Package generator abstracts writing specific Go constructs (interfaces, structs etc) into IO.
package generator

import (
	"bytes"
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

// A Generator is an interface that provides contract for generator struct and is being used instead of a pointer.
type Generator interface {
	WriteHeader()
	WritePackage(pkgName string)
	WriteImports()
	WriteLineBreak(r int)
	WriteInterface(name string, fn ...Func)
	WritePublicStruct(s Struct, usePointers bool)
	WritePrivateStruct(s Struct)
	WriteEnum(e Enum)
	WriteConst(c Const)
	WriteClientConstructor(clientName string)
	WriteInterfaceImplementation(clientName string, f Func)
	WriteGraphqlErrorStructs(usePointers bool)
	Generate() []byte
}

// generator is a private struct that can be created with New function.
type generator struct {
	stream   *bytes.Buffer      // IO to write all code to and read from
	template *template.Template // Predefined template defined in package templates
}

// New return instance of generator.
// rootLoc is relative location of project root (grafik/).
func New(rootLoc string) Generator {
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"sentenceCase": common.SentenceCase,
		"camelCase":    common.SnakeCaseToCamelCase,
	}
	tmpl, err := template.New("codeTemplate").Funcs(funcMap).ParseGlob(filepath.Join(rootLoc, "templates/*.tmpl"))
	if err != nil {
		panic(err)
	}
	return &generator{
		stream:   &bytes.Buffer{},
		template: tmpl.Funcs(funcMap),
	}
}

// WriteHeader writes top level comment (grafik header).
func (g *generator) WriteHeader() {
	g.write(Header)
}

// WritePackage writes name of the package.
func (g *generator) WritePackage(pkgName string) {
	g.write(fmt.Sprintf("package %s", pkgName))
}

// WriteImports writes list of all required imports.
func (g *generator) WriteImports() {
	err := g.template.ExecuteTemplate(g.stream, "imports.tmpl", make(map[string]interface{}))
	if err != nil {
		panic(err)
	}
}

// WriteLineBreak writes number of line breaks based on provided number (r).
func (g *generator) WriteLineBreak(r int) {
	g.write(strings.Repeat("\n", r))
}

// WriteInterface writes interface of provided name and functions (fn).
func (g *generator) WriteInterface(name string, fn ...Func) {
	config := map[string]interface{}{
		"InterfaceName": name,
		"Functions":     fn,
	}
	err := g.template.ExecuteTemplate(g.stream, "interface.tmpl", config)
	if err != nil {
		panic(err)
	}
}

// WritePublicStruct writes struct with capitalized name, fields and json tags based on generator.Struct.
func (g *generator) WritePublicStruct(s Struct, usePointers bool) {
	config := map[string]interface{}{
		"Struct":      s,
		"Public":      true,
		"UsePointers": usePointers,
	}
	err := g.template.ExecuteTemplate(g.stream, "struct.tmpl", config)
	if err != nil {
		panic(err)
	}
}

// WritePrivateStruct writes struct with lowercase name, fields and no json tags based on generator.Struct.
func (g *generator) WritePrivateStruct(s Struct) {
	config := map[string]interface{}{
		"Struct": s,
		"Public": false,
	}
	err := g.template.ExecuteTemplate(g.stream, "struct.tmpl", config)
	if err != nil {
		panic(err)
	}
}

// WriteEnum writes Enum based on generator.Enum.
func (g *generator) WriteEnum(e Enum) {
	err := g.template.ExecuteTemplate(g.stream, "enum.tmpl", e)
	if err != nil {
		panic(err)
	}
}

// WriteConst writes Const based on generator.Const.
func (g *generator) WriteConst(c Const) {
	err := g.template.ExecuteTemplate(g.stream, "const.tmpl", c)
	if err != nil {
		panic(err)
	}
}

// WriteClientConstructor writes New function that serves as a constructor for grafik client.
func (g *generator) WriteClientConstructor(clientName string) {
	err := g.template.ExecuteTemplate(g.stream, "constructor.tmpl", clientName)
	if err != nil {
		panic(err)
	}
}

// WriteInterfaceImplementation writes implementation for earlier defined interface as function with receiver.
func (g *generator) WriteInterfaceImplementation(clientName string, f Func) {
	config := map[string]interface{}{
		"ClientName": clientName,
		"Func":       f,
	}
	err := g.template.ExecuteTemplate(g.stream, "interface_impl.tmpl", config)
	if err != nil {
		panic(err)
	}
}

// WriteGraphqlErrorStructs writes predefined GraphQL error structs.
func (g *generator) WriteGraphqlErrorStructs(usePointers bool) {
	config := map[string]interface{}{
		"UsePointers":            usePointers,
		"GraphQLErrorStructName": GraphQLErrorStructName,
	}
	err := g.template.ExecuteTemplate(g.stream, "graphql_error.tmpl", config)
	if err != nil {
		panic(err)
	}
}

// Generate formats the generated code and returns it as a slice of bytes.
func (g *generator) Generate() []byte {
	b := make([]byte, g.stream.Len())
	_, err := g.stream.Read(b)
	if err != nil {
		panic(err)
	}
	b, err = format.Source(b)
	if err != nil {
		panic(err)
	}
	return b
}

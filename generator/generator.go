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
	WriteImports() error
	WriteLineBreak(r int)
	WriteInterface(name string, fn ...Func) error
	WritePublicStruct(s Struct, usePointers bool) error
	WritePrivateStruct(s Struct) error
	WriteEnum(e Enum) error
	WriteConst(c Const) error
	WriteClientConstructor(clientName string) error
	WriteInterfaceImplementation(clientName string, f Func) error
	WriteGraphqlErrorStructs(usePointers bool) error
	Generate() ([]byte, error)
}

// generator is a private struct that can be created with New function.
type generator struct {
	stream   *bytes.Buffer      // IO to write all code to and read from
	template *template.Template // Predefined template defined in package templates
}

// New return instance of generator.
// rootLoc is relative location of project root (grafik/).
func New(rootLoc string) (Generator, error) {
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"sentenceCase": common.SentenceCase,
		"camelCase":    common.SnakeCaseToCamelCase,
	}
	tmpl, err := template.New("codeTemplate").Funcs(funcMap).ParseGlob(filepath.Join(rootLoc, "templates/*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates. Cause: %w", err)
	}
	return &generator{
		stream:   &bytes.Buffer{},
		template: tmpl.Funcs(funcMap),
	}, nil
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
func (g *generator) WriteImports() error {
	err := g.template.ExecuteTemplate(g.stream, "imports.tmpl", make(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("failed to execute 'imports' template. Cause: %w", err)
	}
	return nil
}

// WriteLineBreak writes number of line breaks based on provided number (r).
func (g *generator) WriteLineBreak(r int) {
	g.write(strings.Repeat("\n", r))
}

// WriteInterface writes interface of provided name and functions (fn).
func (g *generator) WriteInterface(name string, fn ...Func) error {
	config := map[string]interface{}{
		"InterfaceName": name,
		"Functions":     fn,
	}
	err := g.template.ExecuteTemplate(g.stream, "interface.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to execute 'interface' template. Cause: %w", err)
	}
	return nil
}

// WritePublicStruct writes struct with capitalized name, fields and json tags based on generator.Struct.
func (g *generator) WritePublicStruct(s Struct, usePointers bool) error {
	config := map[string]interface{}{
		"Struct":      s,
		"Public":      true,
		"UsePointers": usePointers,
	}
	err := g.template.ExecuteTemplate(g.stream, "struct.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to execute 'struct' template. Cause: %w", err)
	}
	return nil
}

// WritePrivateStruct writes struct with lowercase name, fields and no json tags based on generator.Struct.
func (g *generator) WritePrivateStruct(s Struct) error {
	config := map[string]interface{}{
		"Struct": s,
		"Public": false,
	}
	err := g.template.ExecuteTemplate(g.stream, "struct.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to execute 'struct' template. Cause: %w", err)
	}
	return nil
}

// WriteEnum writes Enum based on generator.Enum.
func (g *generator) WriteEnum(e Enum) error {
	err := g.template.ExecuteTemplate(g.stream, "enum.tmpl", e)
	if err != nil {
		return fmt.Errorf("failed to execute 'enum' template. Cause: %w", err)
	}
	return nil
}

// WriteConst writes Const based on generator.Const.
func (g *generator) WriteConst(c Const) error {
	err := g.template.ExecuteTemplate(g.stream, "const.tmpl", c)
	if err != nil {
		return fmt.Errorf("failed to execute 'const' template. Cause: %w", err)
	}
	return nil
}

// WriteClientConstructor writes New function that serves as a constructor for grafik client.
func (g *generator) WriteClientConstructor(clientName string) error {
	err := g.template.ExecuteTemplate(g.stream, "constructor.tmpl", clientName)
	if err != nil {
		return fmt.Errorf("failed to execute 'constructor' template. Cause: %w", err)
	}
	return nil
}

// WriteInterfaceImplementation writes implementation for earlier defined interface as function with receiver.
func (g *generator) WriteInterfaceImplementation(clientName string, f Func) error {
	config := map[string]interface{}{
		"ClientName": clientName,
		"Func":       f,
	}
	err := g.template.ExecuteTemplate(g.stream, "interface_impl.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to execute 'interface_impl' template. Cause: %w", err)
	}
	return nil
}

// WriteGraphqlErrorStructs writes predefined GraphQL error structs.
func (g *generator) WriteGraphqlErrorStructs(usePointers bool) error {
	config := map[string]interface{}{
		"UsePointers":            usePointers,
		"GraphQLErrorStructName": GraphQLErrorStructName,
	}
	err := g.template.ExecuteTemplate(g.stream, "graphql_error.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to execute 'graphql_error' template. Cause: %w", err)
	}
	return nil
}

// Generate formats the generated code and returns it as a slice of bytes.
func (g *generator) Generate() ([]byte, error) {
	b := make([]byte, g.stream.Len())
	_, err := g.stream.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream content from generator. Cause: %w", err)
	}
	b, err = format.Source(b)
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code. Cause: %w", err)
	}
	return b, nil
}

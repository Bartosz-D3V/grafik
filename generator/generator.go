// Package generator abstracts writing specific Go constructs (interfaces, structs etc) into IO.
package generator

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"go/format"
	"go/parser"
	"go/token"
	"io"
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
	Generate() io.WriterTo
}

//go:embed templates/*
var content embed.FS

type generatorWriter interface {
	io.StringWriter
	io.Writer
}

// generator is a private struct that can be created with New function.
type generator struct {
	// IO to write all code to and read from.
	stream generatorWriter
	// Predefined template defined in package templates.
	template *template.Template
}

// New return instance of generator.
func New() Generator {
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"sentenceCase": common.SentenceCase,
		"camelCase":    common.SnakeCaseToCamelCase,
	}

	tmpl, err := template.New("codeTemplate").Funcs(funcMap).ParseFS(content, "**/*.tmpl")
	if err != nil {
		panic(fmt.Errorf("failed to parse templates. Cause: %w", err))
	}
	return &generator{
		stream:   &bytes.Buffer{},
		template: tmpl.Funcs(funcMap),
	}
}

// WriteHeader writes top level comment (grafik header).
func (g *generator) WriteHeader() {
	_, err := g.stream.WriteString(Header)
	if err != nil {
		panic(fmt.Errorf("failed to write header. Cause: %w", err))
	}
}

// WritePackage writes name of the package.
func (g *generator) WritePackage(pkgName string) {
	_, err := g.stream.WriteString(fmt.Sprintf("package %s", pkgName))
	if err != nil {
		panic(fmt.Errorf("failed to write package name. Cause: %w", err))
	}
}

// WriteImports writes list of all required imports.
func (g *generator) WriteImports() {
	err := g.template.ExecuteTemplate(g.stream, "imports.tmpl", make(map[string]interface{}))
	if err != nil {
		panic(fmt.Errorf("failed to execute 'imports' template. Cause: %w", err))
	}
}

// WriteLineBreak writes number of line breaks based on provided number (r).
func (g *generator) WriteLineBreak(r int) {
	_, err := g.stream.WriteString(strings.Repeat("\n", r))
	if err != nil {
		panic(fmt.Errorf("failed to write line break. Cause: %w", err))
	}
}

// WriteInterface writes interface of provided name and functions (fn).
func (g *generator) WriteInterface(name string, fn ...Func) {
	config := map[string]interface{}{
		"InterfaceName": name,
		"Functions":     fn,
	}
	err := g.template.ExecuteTemplate(g.stream, "interface.tmpl", config)
	if err != nil {
		panic(fmt.Errorf("failed to execute 'interface' template. Cause: %w", err))
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
		panic(fmt.Errorf("failed to execute 'struct' template. Cause: %w", err))
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
		panic(fmt.Errorf("failed to execute 'struct' template. Cause: %w", err))
	}
}

// WriteEnum writes Enum based on generator.Enum.
func (g *generator) WriteEnum(e Enum) {
	err := g.template.ExecuteTemplate(g.stream, "enum.tmpl", e)
	if err != nil {
		panic(fmt.Errorf("failed to execute 'enum' template. Cause: %w", err))
	}
}

// WriteConst writes Const based on generator.Const.
func (g *generator) WriteConst(c Const) {
	err := g.template.ExecuteTemplate(g.stream, "const.tmpl", c)
	if err != nil {
		panic(fmt.Errorf("failed to execute 'const' template. Cause: %w", err))
	}
}

// WriteClientConstructor writes New function that serves as a constructor for grafik client.
func (g *generator) WriteClientConstructor(clientName string) {
	err := g.template.ExecuteTemplate(g.stream, "constructor.tmpl", clientName)
	if err != nil {
		panic(fmt.Errorf("failed to execute 'constructor' template. Cause: %w", err))
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
		panic(fmt.Errorf("failed to execute 'interface_impl' template. Cause: %w", err))
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
		panic(fmt.Errorf("failed to execute 'graphql_error' template. Cause: %w", err))
	}
}

// Generate formats the generated code and returns it as a WriterTo interface.
func (g *generator) Generate() io.WriterTo {
	writer := &bytes.Buffer{}
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", g.stream, parser.ParseComments)
	if err != nil {
		panic(fmt.Errorf("failed to parse generated Go code. Cause: %w", err))
	}

	err = format.Node(writer, fSet, f)
	if err != nil {
		panic(fmt.Errorf("failed to gofmt generated Go code. Cause: %w", err))
	}
	return writer
}

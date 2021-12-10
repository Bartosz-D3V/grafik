package generator

import (
	"bytes"
	"fmt"
	"github.com/Bartosz-D3V/ggrafik/common"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

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

type generator struct {
	stream   *bytes.Buffer
	template *template.Template
}

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

func (g *generator) WriteHeader() {
	g.write(Header)
}

func (g *generator) WritePackage(pkgName string) {
	g.write(fmt.Sprintf("package %s", pkgName))
}

func (g *generator) WriteImports() {
	err := g.template.ExecuteTemplate(g.stream, "imports.tmpl", make(map[string]interface{}))
	if err != nil {
		panic(err)
	}
}

func (g *generator) WriteLineBreak(r int) {
	g.write(strings.Repeat("\n", r))
}

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

func (g *generator) WriteEnum(e Enum) {
	err := g.template.ExecuteTemplate(g.stream, "enum.tmpl", e)
	if err != nil {
		panic(err)
	}
}

func (g *generator) WriteConst(c Const) {
	err := g.template.ExecuteTemplate(g.stream, "const.tmpl", c)
	if err != nil {
		panic(err)
	}
}

func (g *generator) WriteClientConstructor(clientName string) {
	err := g.template.ExecuteTemplate(g.stream, "constructor.tmpl", clientName)
	if err != nil {
		panic(err)
	}
}

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

func (g *generator) WriteGraphqlErrorStructs(usePointers bool) {
	config := map[string]interface{}{
		"UsePointers": usePointers,
	}
	err := g.template.ExecuteTemplate(g.stream, "graphql_error.tmpl", config)
	if err != nil {
		panic(err)
	}
}

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

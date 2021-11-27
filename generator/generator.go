package generator

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator interface {
	WriteHeader()
	WriteLineBreak(r int)
	WriteInterface(name string, fn ...Func)
	WritePublicStruct(s Struct)
	WritePrivateStruct(s Struct)
	WriteEnum(e Enum)
	WriteConst(c Const)
	WriteClientConstructor(clientName string)
	Generate() []byte
}

type generator struct {
	stream   *bytes.Buffer
	template *template.Template
}

func New(fptr string) Generator {
	tmpl, err := template.ParseGlob(filepath.Join(fptr, "templates/*.tmpl"))
	if err != nil {
		panic(err)
	}
	return &generator{
		stream:   &bytes.Buffer{},
		template: tmpl,
	}
}

func (g *generator) WriteHeader() {
	g.write(Header)
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

func (g *generator) WritePublicStruct(s Struct) {
	config := map[string]interface{}{
		"Struct": s,
		"Public": true,
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

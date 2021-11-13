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
	WriteStruct(s Struct)
	Generate() []byte
}

type generator struct {
	stream   bytes.Buffer
	ident    int
	template *template.Template
}

func New(fptr string) Generator {
	tmpl, err := template.ParseGlob(filepath.Join(fptr, "templates/*.tmpl"))
	if err != nil {
		panic(err)
	}
	return &generator{
		stream:   bytes.Buffer{},
		ident:    0,
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
	err := g.template.ExecuteTemplate(&g.stream, "interface.tmpl", config)
	if err != nil {
		panic(err)
	}
}

func (g *generator) WriteStruct(s Struct) {
	err := g.template.ExecuteTemplate(&g.stream, "struct.tmpl", s)
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

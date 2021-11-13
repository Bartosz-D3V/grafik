package generator

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
	stream   bytes.Buffer
	ident    int
	template *template.Template
}

func New(fptr string) *Generator {
	tmpl, err := template.ParseGlob(filepath.Join(fptr, "templates/*.tmpl"))
	if err != nil {
		panic(err)
	}
	return &Generator{
		stream:   bytes.Buffer{},
		ident:    0,
		template: tmpl,
	}
}

func (g *Generator) WriteHeader() {
	g.write(Header)
}

func (g *Generator) WriteLineBreaks(r int) {
	g.write(strings.Repeat("\n", r))
}

func (g *Generator) WriteInterface(name string, fn ...Func) {
	config := map[string]interface{}{
		"InterfaceName": name,
		"Functions":     fn,
	}
	err := g.template.ExecuteTemplate(&g.stream, "interface.tmpl", config)
	if err != nil {
		return
	}
}

func (g *Generator) Generate() []byte {
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

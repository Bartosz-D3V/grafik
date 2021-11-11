package generator

import (
	"bytes"
	"text/template"
)

type Generator struct {
	stream   bytes.Buffer
	ident    int
	template *template.Template
}

func New() *Generator {
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic("pa")
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
	return g.stream.Bytes()
}

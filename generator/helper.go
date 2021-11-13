package generator

import (
	"fmt"
	"strings"
)

func (g *generator) write(s string) {
	id := strings.Repeat("\t", g.ident)
	g.stream.WriteString(id)
	g.stream.WriteString(s)
}

func (g *generator) sWrite(s string, args ...interface{}) {
	f := fmt.Sprintf(s, args...)
	id := strings.Repeat("\t", g.ident)
	g.stream.WriteString(id)
	g.stream.WriteString(f)
}

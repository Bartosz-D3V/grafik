package generator

import (
	"fmt"
)

func (g *generator) write(s string) {
	g.stream.WriteString(s)
}

func (g *generator) sWrite(s string, args ...interface{}) {
	f := fmt.Sprintf(s, args...)
	g.stream.WriteString(f)
}

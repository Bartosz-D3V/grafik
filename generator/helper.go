package generator

func (g *generator) write(s string) {
	g.stream.WriteString(s)
}

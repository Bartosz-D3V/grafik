package evaluator

import (
	"fmt"
	"github.com/Bartosz-D3V/ggrafik/test"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"path"
	"testing"
)

func TestEvaluator_Generate_FlatStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/simple-definition.graphql")
	e := New(pd, schema)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	file(id string) File
	files()
}

type Mutation interface {
	renameFile(id string, name string) File
	deleteFile(id string) File
}

type File struct {
	name string
}
`))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_NestedStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/nested-type-definition.graphql")
	e := New(pd, schema)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	getHero() Character
}

type Character struct {
	name      string
	homeWorld Planet
	species   Species
}

type Planet struct {
	name     string
	climate  string
	location Location
}

type Species struct {
	name     string
	lifespan int
	origin   Planet
}

type Location struct {
	posX int
	poxY int
}
`))

	assert.Equal(t, expOut, out)
}

func loadSchema(t *testing.T, pd string, schemaName string) *ast.Schema {
	schemaLoc := path.Join(pd, schemaName)
	file, err := ioutil.ReadFile(schemaLoc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(file),
		Name:  "simple-definition.graphql",
	})
	assert.Nil(t, err)
	assert.NotNil(t, schema)

	return schema
}

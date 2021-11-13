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

func TestEvaluator_Generate(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd)
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

func loadSchema(t *testing.T, pd string) *ast.Schema {
	schemaLoc := path.Join(pd, "test/simple-definition.graphql")
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

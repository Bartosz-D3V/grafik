package visitor

import (
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"path"
	"testing"
)

func TestVisitor_IntrospectTypes_RepeatedFields(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "repeated_fields/schema.graphql")
	query := loadQuery(t, schema, "repeated_fields/query.graphql")

	v := New(schema, query)
	types := v.IntrospectTypes()

	assert.Equal(t, 10, len(types["Character"]))
	assert.Equal(t, 4, len(types["Friend"]))
}

func loadSchema(t *testing.T, loc string) *ast.Schema {
	file, err := ioutil.ReadFile(loc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadSchema(&ast.Source{
		Input: string(file),
		Name:  path.Base(loc),
	})
}

func loadQuery(t *testing.T, schema *ast.Schema, loc string) *ast.QueryDocument {
	file, err := ioutil.ReadFile(loc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadQuery(schema, string(file))
}

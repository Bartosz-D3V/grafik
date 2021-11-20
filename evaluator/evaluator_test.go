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
	schema := loadSchema(t, pd, "test/simple_type/simple_type.graphql")
	query := loadQuery(t, pd, schema, "test/simple_type/simple_type_query.graphql")
	e := New(pd, schema, query)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	file(id string) File
	files() []File
}

type Mutation interface {
	renameFile(id string, name string) File
	deleteFile(id string) File
}

type File struct {
	name string
}

const getFile = %cquery getFile {
    file(id: "123-ABC") {
        name
    }
}%c
`, '`', '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_ArrayStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/array/array.graphql")
	query := loadQuery(t, pd, schema, "test/array/array_query.graphql")
	e := New(pd, schema, query)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	getBook() Book
}

type Book struct {
	name string
	tags []Tag
}

type Tag struct {
	name string
}

const getBookTags = %cquery getBookTags {
    getBook {
        tags {
            name
        }
        name
    }
}%c
`, '`', '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_NestedStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/nested_type/nested_type.graphql")
	query := loadQuery(t, pd, schema, "test/nested_type/nested_type_query.graphql")
	e := New(pd, schema, query)

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

const getHero = %cquery getHero {
    getHero {
        name
    }
}%c
`, '`', '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Enum(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/enum/enum.graphql")
	query := loadQuery(t, pd, schema, "test/enum/enum_query.graphql")
	e := New(pd, schema, query)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	getDepartment() Department
}

type DepartmentName string

const (
	IT      DepartmentName = "IT"
	SALES   DepartmentName = "SALES"
	HR      DepartmentName = "HR"
	SUPPORT DepartmentName = "SUPPORT"
)

type Department struct {
	name DepartmentName
}

const getDepartment = %cquery getDepartment {
    getDepartment {
        name
    }
}%c
`, '`', '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Input(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/input/input.graphql")
	query := loadQuery(t, pd, schema, "test/input/input_query.graphql")
	e := New(pd, schema, query)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Query interface {
	all(company Company) string
}

type Company struct {
	code int
	eq   string
}

const getCompanyWithCode123 = %cquery getCompanyWithCode123 {
    all(company: {code: 123})
}%c
`, '`', '`'))

	assert.Equal(t, expOut, out)
}

func loadSchema(t *testing.T, pd string, schemaName string) *ast.Schema {
	schemaLoc := path.Join(pd, schemaName)
	file, err := ioutil.ReadFile(schemaLoc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadSchema(&ast.Source{
		Input: string(file),
		Name:  path.Base(schemaName),
	})
}

func loadQuery(t *testing.T, pd string, schema *ast.Schema, queryName string) *ast.QueryDocument {
	queryLoc := path.Join(pd, queryName)
	file, err := ioutil.ReadFile(queryLoc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadQuery(schema, string(file))
}

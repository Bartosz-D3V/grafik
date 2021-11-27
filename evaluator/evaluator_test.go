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
	e := New(pd, schema, query, "SimpleType")

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type File struct {
	Name string %[1]cjson:"name"%[1]c
}

const getFile = %[1]cquery getFile {
    file(id: "123-ABC") {
        name
    }
}%[1]c

type SimpleTypeGraphql interface {
	GetFile() (*http.Response, error)
}

type SimpleTypeGraphql struct {
	ctrl graphqlClient.Client
}

func New(endpoint string, client *http.Client) *SimpleTypeGraphql {
	return &SimpleTypeGraphql{
		ctrl: graphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_ArrayStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/array/array.graphql")
	query := loadQuery(t, pd, schema, "test/array/array_query.graphql")
	e := New(pd, schema, query, "")

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Book struct {
	Name string %[1]cjson:"name"%[1]c
	Tags []Tag  %[1]cjson:"tags"%[1]c
}

type Tag struct {
	Name string %[1]cjson:"name"%[1]c
}

const getBookTags = %[1]cquery getBookTags {
    getBook {
        tags {
            name
        }
        name
    }
}%[1]c

type ArrayGraphql interface {
	GetBookTags() (*http.Response, error)
}

type arrayGraphql struct {
	ctrl graphqlClient.Client
}

func New(endpoint string, client *http.Client) *arrayGraphql {
	return &arrayGraphql{
		ctrl: graphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_NestedStructure(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/nested_type/nested_type.graphql")
	query := loadQuery(t, pd, schema, "test/nested_type/nested_type_query.graphql")
	e := New(pd, schema, query, "NestedType")

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Character struct {
	Name      string  %[1]cjson:"name"%[1]c
	HomeWorld Planet  %[1]cjson:"homeWorld"%[1]c
	Species   Species %[1]cjson:"species"%[1]c
}

type Planet struct {
	Name     string   %[1]cjson:"name"%[1]c
	Climate  string   %[1]cjson:"climate"%[1]c
	Location Location %[1]cjson:"location"%[1]c
}

type Location struct {
	PosX int %[1]cjson:"posX"%[1]c
	PoxY int %[1]cjson:"poxY"%[1]c
}

type Species struct {
	Name     string %[1]cjson:"name"%[1]c
	Lifespan int    %[1]cjson:"lifespan"%[1]c
	Origin   Planet %[1]cjson:"origin"%[1]c
}

type CharacterSelector struct {
	IdSelector IdSelector %[1]cjson:"idSelector"%[1]c
}

type IdSelector struct {
	Id string %[1]cjson:"id"%[1]c
}

const getHero = %[1]cquery getHero {
    getHero(characterSelector: {idSelector: {id: "123ABC"}}) {
        homeWorld {
            location {
                posX
            }
        }
    }
}%[1]c

type NestedTypeGraphql interface {
	GetHero() (*http.Response, error)
}

type NestedTypeGraphql struct {
	ctrl graphqlClient.Client
}

func New(endpoint string, client *http.Client) *NestedTypeGraphql {
	return &NestedTypeGraphql{
		ctrl: graphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Enum(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/enum/enum.graphql")
	query := loadQuery(t, pd, schema, "test/enum/enum_query.graphql")
	e := New(pd, schema, query, "")

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type DepartmentName string

const (
	IT      DepartmentName = "IT"
	SALES   DepartmentName = "SALES"
	HR      DepartmentName = "HR"
	SUPPORT DepartmentName = "SUPPORT"
)

type Department struct {
	Name DepartmentName %[1]cjson:"name"%[1]c
}

const getDepartment = %[1]cquery getDepartment {
    getDepartment {
        name
    }
}%[1]c

type EnumGraphql interface {
	GetDepartment() (*http.Response, error)
}

type enumGraphql struct {
	ctrl graphqlClient.Client
}

func New(endpoint string, client *http.Client) *enumGraphql {
	return &enumGraphql{
		ctrl: graphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Input(t *testing.T) {
	pd := test.GetParentDir(t)
	schema := loadSchema(t, pd, "test/input/input.graphql")
	query := loadQuery(t, pd, schema, "test/input/input_query.graphql")
	e := New(pd, schema, query, "")

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

type Company struct {
	Code int    %[1]cjson:"code"%[1]c
	Eq   string %[1]cjson:"eq"%[1]c
}

const getCompanyWithCode123 = %[1]cquery getCompanyWithCode123 {
    all(company: {code: 123})
}%[1]c

type InputGraphql interface {
	GetCompanyWithCode123() (*http.Response, error)
}

type inputGraphql struct {
	ctrl graphqlClient.Client
}

func New(endpoint string, client *http.Client) *inputGraphql {
	return &inputGraphql{
		ctrl: graphqlClient.New(endpoint, client),
	}
}
`, '`'))

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

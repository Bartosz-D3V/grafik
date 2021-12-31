package generator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Bartosz-D3V/grafik/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
)

func TestGenerator_WriteHeader(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WriteHeader()
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WritePackage("test")
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
// Generated with grafik. DO NOT EDIT

package test
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteHeader_Error(t *testing.T) {
	t.Parallel()

	g := generator{stream: faultyWriter{}}

	err := g.WriteHeader()
	assert.EqualError(t, err, "failed to write header. Cause: unit test: Failed to write a string")
}

func TestGenerator_WritePackage_Error(t *testing.T) {
	t.Parallel()

	g := generator{stream: faultyWriter{}}

	err := g.WritePackage("")
	assert.EqualError(t, err, "failed to write package name. Cause: unit test: Failed to write a string")
}

func TestGenerator_WriteLineBreak_Error(t *testing.T) {
	t.Parallel()

	g := generator{stream: faultyWriter{}}

	err := g.WriteLineBreak(1)
	assert.EqualError(t, err, "failed to write line break. Cause: unit test: Failed to write a string")
}

func TestGenerator_WriteImports(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteImports()
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
package test

import (
    "context"
    GraphqlClient "github.com/Bartosz-D3V/grafik/client"
    "net/http"
)
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteImports_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("imports.tmpl").Parse("imports.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteImports()
	assert.EqualError(t, err, "failed to execute 'imports' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteInterface_NoArgWithReturn(t *testing.T) {
	t.Parallel()
	fn := Func{
		Name: "FindBook",
		Args: make([]TypeArg, 0),
		Type: "Book",
	}

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	_ = g.WriteInterface("BookService", fn)

	out := getSourceString(t, g)

	expOut := test.PrepExpCode(t, `
package test

type BookService interface {
	FindBook(ctx context.Context, header *http.Header) Book
}`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_SingleArgWithReturn(t *testing.T) {
	t.Parallel()
	fn := Func{
		Name: "FindBook",
		Args: []TypeArg{{Name: "isbn", Type: "string"}},
		Type: "Book",
	}

	g, err := New()
	assert.NoError(t, err)

	err = g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteInterface("BookService", fn)
	assert.NoError(t, err)

	out := getSourceString(t, g)

	expOut := test.PrepExpCode(t, `
package test

type BookService interface {
	FindBook(ctx context.Context, isbn string, header *http.Header) Book
}`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_MultiArgWithReturn(t *testing.T) {
	t.Parallel()
	fn := Func{
		Name: "FindEmployee",
		Args: []TypeArg{
			{Name: "name", Type: "string"},
			{Name: "department", Type: "string"},
			{Name: "age", Type: "int"},
		},
		Type: "Employee",
	}

	g, err := New()
	assert.NoError(t, err)

	err = g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteInterface("EmployeeService", fn)
	assert.NoError(t, err)

	out := getSourceString(t, g)

	expOut := test.PrepExpCode(t, `
package test

type EmployeeService interface {
	FindEmployee(ctx context.Context, name string, department string, age int, header *http.Header) Employee
}`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_MultiMethods(t *testing.T) {
	t.Parallel()
	fn1 := Func{
		Name: "FindBook",
		Args: make([]TypeArg, 0),
		Type: "Book",
	}
	fn2 := Func{
		Name: "FindEmployee",
		Args: []TypeArg{
			{Name: "name", Type: "string"},
			{Name: "department", Type: "string"},
			{Name: "age", Type: "int"},
		},
		Type: "Employee",
	}

	g, err := New()
	assert.NoError(t, err)

	err = g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteInterface("BookService", fn1, fn2)
	assert.NoError(t, err)

	out := getSourceString(t, g)

	expOut := test.PrepExpCode(t, `
package test

type BookService interface {
	FindBook(ctx context.Context, header *http.Header) Book
	FindEmployee(ctx context.Context, name string, department string, age int, header *http.Header) Employee
}`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("interface.tmpl").Parse("interface.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteInterface("")
	assert.EqualError(t, err, "failed to execute 'interface' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WritePublicStruct(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	s := Struct{
		Name: "Person",
		Fields: []TypeArg{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "int",
			},
		},
	}
	err = g.WritePublicStruct(s, false)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
package test

type Person struct {
	Name string %[1]cjson:"name"%[1]c
	Age  int    %[1]cjson:"age"%[1]c
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WritePublicStruct_WithPointers(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	s := Struct{
		Name: "Person",
		Fields: []TypeArg{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "int",
			},
		},
	}
	err = g.WritePublicStruct(s, true)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
package test

type Person struct {
	Name *string %[1]cjson:"name"%[1]c
	Age  *int    %[1]cjson:"age"%[1]c
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WritePublicStruct_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("struct.tmpl").Parse("struct.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WritePublicStruct(Struct{}, false)
	assert.EqualError(t, err, "failed to execute 'struct' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WritePrivateStruct(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	s := Struct{
		Name: "Person",
		Fields: []TypeArg{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "int",
			},
		},
	}
	err = g.WritePrivateStruct(s)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
package test

type person struct {
	name string
	age  int
}
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WritePrivateStruct_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("struct.tmpl").Parse("struct.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WritePrivateStruct(Struct{})
	assert.EqualError(t, err, "failed to execute 'struct' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteEnum(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	e := Enum{
		Name:   "Planet",
		Fields: []string{"NEPTUNE", "MARS", "SATURN"},
	}
	err = g.WriteEnum(e)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
package test

type Planet string

const (
	NEPTUNE Planet = "NEPTUNE"
	MARS    Planet = "MARS"
	SATURN  Planet = "SATURN"
)
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteEnum_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("enum.tmpl").Parse("enum.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteEnum(Enum{})
	assert.EqualError(t, err, "failed to execute 'enum' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteConst(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	c1 := Const{
		Name: "pi",
		Val:  3.16,
	}
	c2 := Const{
		Name: "encoding",
		Val:  "UTF-8",
	}
	err = g.WriteConst(c1)
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteConst(c2)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
package test

const pi = 3.16

const encoding = %[1]cUTF-8%[1]c
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteConst_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("const.tmpl").Parse("const.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteConst(Const{})
	assert.EqualError(t, err, "failed to execute 'const' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteClientConstructor(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteClientConstructor("ApiClient")
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
package test

func New(endpoint string, client *http.Client) ApiClient {
    return &apiClient {
        ctrl: GraphqlClient.New(endpoint, client),
    }
}
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteClientConstructor_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("constructor.tmpl").Parse("constructor.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteClientConstructor("")
	assert.EqualError(t, err, "failed to execute 'constructor' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteInterfaceImplementation(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	f := Func{
		Name: "countResults",
		Args: []TypeArg{{
			Name: "condition",
			Type: "string",
		}},
		Type:         "int",
		WrapperTypes: nil,
	}
	err = g.WriteInterfaceImplementation("apiClient", f)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, `
package test

func (c *apiClient) CountResults(ctx context.Context, condition string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["condition"] = condition

	return c.ctrl.Execute(ctx, countResults, params, header)
}
`)

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterfaceImplementation_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("interface_impl.tmpl").Parse("interface_impl.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteInterfaceImplementation("", Func{})
	assert.EqualError(t, err, "failed to execute 'interface_impl' template. Cause: unit test: Failed to write a slice of bytes")
}

func TestGenerator_WriteGraphqlErrorStructs(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteGraphqlErrorStructs(false)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
package test

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteGraphqlErrorStructs_WithPointers(t *testing.T) {
	t.Parallel()

	g, _ := New()

	err := g.WritePackage("test")
	assert.NoError(t, err)

	err = g.WriteLineBreak(2)
	assert.NoError(t, err)

	err = g.WriteGraphqlErrorStructs(true)
	assert.NoError(t, err)

	out := getSourceString(t, g)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
package test

type GraphQLError struct {
	Message    *string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation  %[1]cjson:"locations"%[1]c
	Extensions *GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   *int %[1]cjson:"line"%[1]c
	Column *int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code *string %[1]cjson:"code"%[1]c
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func getSourceString(t *testing.T, g Generator) string {
	fileContent, err := g.Generate()
	assert.NoError(t, err)
	src := &bytes.Buffer{}
	_, err = fileContent.WriteTo(src)
	assert.NoError(t, err)

	return src.String()
}

func TestGenerator_WriteGraphqlErrorStructs_Error(t *testing.T) {
	t.Parallel()

	pTmpl, _ := template.New("graphql_error.tmpl").Parse("graphql_error.tmpl")
	g := generator{
		stream:   faultyWriter{},
		template: pTmpl,
	}

	err := g.WriteGraphqlErrorStructs(false)
	assert.EqualError(t, err, "failed to execute 'graphql_error' template. Cause: unit test: Failed to write a slice of bytes")
}

type faultyWriter struct{}

func (w faultyWriter) WriteString(string) (n int, err error) {
	return 0, errors.New("unit test: Failed to write a string")
}

func (w faultyWriter) Write([]byte) (n int, err error) {
	return 0, errors.New("unit test: Failed to write a slice of bytes")
}

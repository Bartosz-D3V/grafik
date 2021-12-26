package generator

import (
	"github.com/Bartosz-D3V/grafik/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerator_WriteInterface_NoArgWithReturn(t *testing.T) {
	t.Parallel()
	fn := Func{
		Name: "FindBook",
		Args: make([]TypeArg, 0),
		Type: "Book",
	}

	g, _ := New("../")
	_ = g.WriteInterface("BookService", fn)

	fileContent, err := g.Generate()
	assert.NoError(t, err)
	out := string(fileContent)

	expOut := test.PrepExpCode(t, `
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

	g, err := New("../")
	assert.NoError(t, err)

	err = g.WriteInterface("BookService", fn)
	assert.NoError(t, err)

	fileContent, err := g.Generate()
	assert.NoError(t, err)
	out := string(fileContent)

	expOut := test.PrepExpCode(t, `
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

	g, err := New("../")
	assert.NoError(t, err)

	err = g.WriteInterface("EmployeeService", fn)
	assert.NoError(t, err)

	fileContent, err := g.Generate()
	assert.NoError(t, err)
	out := string(fileContent)

	expOut := test.PrepExpCode(t, `
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

	g, err := New("../")
	assert.NoError(t, err)

	err = g.WriteInterface("BookService", fn1, fn2)
	assert.NoError(t, err)

	fileContent, err := g.Generate()
	assert.NoError(t, err)
	out := string(fileContent)

	expOut := test.PrepExpCode(t, `
type BookService interface {
	FindBook(ctx context.Context, header *http.Header) Book
	FindEmployee(ctx context.Context, name string, department string, age int, header *http.Header) Employee
}`)

	assert.Equal(t, expOut, out)
}

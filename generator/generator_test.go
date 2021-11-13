package generator

import (
	"fmt"
	"github.com/Bartosz-D3V/ggrafik/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerator_WriteInterface_NoArgWithReturn(t *testing.T) {
	fn := Func{
		Name: "FindBook",
		Args: make([]FuncArg, 0),
		Type: "Book",
	}
	pd := test.GetParentDir(t)

	g := New(pd)
	g.WriteInterface("BookService", fn)

	out := string(g.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
type BookService interface {
	FindBook() Book
}`))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_SingleArgWithReturn(t *testing.T) {
	fn := Func{
		Name: "FindBook",
		Args: []FuncArg{{Name: "isbn", Type: "string"}},
		Type: "Book",
	}
	pd := test.GetParentDir(t)

	g := New(pd)
	g.WriteInterface("BookService", fn)

	out := string(g.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
type BookService interface {
	FindBook(isbn string) Book
}`))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_MultiArgWithReturn(t *testing.T) {
	fn := Func{
		Name: "FindEmployee",
		Args: []FuncArg{
			{Name: "name", Type: "string"},
			{Name: "department", Type: "string"},
			{Name: "age", Type: "int"},
		},
		Type: "Employee",
	}
	pd := test.GetParentDir(t)

	g := New(pd)
	g.WriteInterface("EmployeeService", fn)

	out := string(g.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
type EmployeeService interface {
	FindEmployee(name string, department string, age int) Employee
}`))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_MultiMethods(t *testing.T) {
	fn1 := Func{
		Name: "FindBook",
		Args: make([]FuncArg, 0),
		Type: "Book",
	}
	fn2 := Func{
		Name: "FindEmployee",
		Args: []FuncArg{
			{Name: "name", Type: "string"},
			{Name: "department", Type: "string"},
			{Name: "age", Type: "int"},
		},
		Type: "Employee",
	}
	pd := test.GetParentDir(t)

	g := New(pd)
	g.WriteInterface("BookService", fn1, fn2)

	out := string(g.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
type BookService interface {
	FindBook() Book
	FindEmployee(name string, department string, age int) Employee
}`))

	assert.Equal(t, expOut, out)
}

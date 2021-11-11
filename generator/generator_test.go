package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGenerator_WriteInterface_NoArgWithReturn(t *testing.T) {
	fn := Func{
		Name: "FindBook",
		Args: make([]FuncArg, 0),
		Type: "Book",
	}
	g := New()
	g.WriteHeader()
	g.WriteInterface("BookService", fn)

	out := string(g.Generate())
	expOut := strings.TrimSpace(fmt.Sprintf(`
%s
type BookService interface {
	FindBook() Book
}`, Header))

	assert.Equal(t, expOut, out)
}

func TestGenerator_WriteInterface_SingleArgWithReturn(t *testing.T) {
	fn := Func{
		Name: "FindBook",
		Args: []FuncArg{{Name: "isbn", Type: "string"}},
		Type: "Book",
	}
	g := New()
	g.WriteHeader()
	g.WriteInterface("BookService", fn)

	out := string(g.Generate())
	expOut := strings.TrimSpace(fmt.Sprintf(`
%s
type BookService interface {
	FindBook(isbn string) Book
}`, Header))

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

	g := New()
	g.WriteHeader()
	g.WriteInterface("EmployeeService", fn)

	out := string(g.Generate())
	expOut := strings.TrimSpace(fmt.Sprintf(`
%s
type EmployeeService interface {
	FindEmployee(name string, department string, age int) Employee
}`, Header))

	assert.Equal(t, expOut, out)
}

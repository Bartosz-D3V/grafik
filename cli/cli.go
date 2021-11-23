package main

import (
	"flag"
	"github.com/Bartosz-D3V/ggrafik/evaluator"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"os"
)

func main() {
	wd, err := os.Getwd()
	fptr := flag.String("fpath", wd, "file path to read from")
	flag.Parse()

	file, err := ioutil.ReadFile("test/simple_type/simple_type.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(file),
		Name:  "simple_type.graphql",
	})

	fileQuery, err := ioutil.ReadFile("test/simple_type/simple_type_query.graphql")
	query, err := gqlparser.LoadQuery(schema, string(fileQuery))

	e := evaluator.New(*fptr, schema, query)

	print(string(e.Generate()))
}

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

	file, err := ioutil.ReadFile("test/countries.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(file),
		Name:  "nested-type-definition.graphql",
	})

	fileQuery, err := ioutil.ReadFile("test/coutries_query.graphql")
	query, err := gqlparser.LoadQuery(schema, string(fileQuery))
	print(query)
	e := evaluator.New(*fptr, schema)
	print(string(e.Generate()))
}

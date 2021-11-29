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
	if err != nil {
		panic(err)
	}
	fptr := flag.String("fpath", wd, "file path to read from")
	flag.Parse()

	clientName := "countries"
	file, err := ioutil.ReadFile("test/countries/countries.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(file),
		Name:  clientName,
	})
	if err != nil {
		panic(err)
	}

	fileQuery, err := ioutil.ReadFile("test/countries/countries_query.graphql")
	if err != nil {
		panic(err)
	}
	query, err := gqlparser.LoadQuery(schema, string(fileQuery))
	if err != nil {
		panic(err)
	}

	e := evaluator.New(*fptr, schema, query, clientName)

	print(string(e.Generate()))
}

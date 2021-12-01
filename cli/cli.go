package main

import (
	"flag"
	"github.com/Bartosz-D3V/ggrafik/evaluator"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"os"
)

type cli struct {
	src        string
	dst        string
	clientName string
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fptr := flag.String("fpath", wd, "file path to read from")
	flag.Parse()

	clientName := "countries"

	pkgName := "client_graphql"

	cli := cli{
		src:        *fptr,
		dst:        "./generated_client/",
		clientName: clientName,
	}

	file, err := ioutil.ReadFile("test/countries/countries.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(file),
		Name:  clientName,
	})
	print(err.Error())

	fileQuery, err := ioutil.ReadFile("test/countries/countries_query.graphql")
	if err != nil {
		panic(err)
	}
	query, err := gqlparser.LoadQuery(schema, string(fileQuery))
	print(err.Error())

	e := evaluator.New(*fptr, schema, query, clientName, pkgName)
	cli.writeFile(e.Generate())
}

package main

import (
	"flag"
	"fmt"
	"github.com/Bartosz-D3V/grafik/evaluator"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"os"
)

type cli struct {
	schemaSource *string
	querySource  *string
	packageName  *string
	clientName   *string
	destination  *string
	usePointers  *bool
}

func main() {
	genCmd := flag.NewFlagSet("gen", flag.ExitOnError)
	genSchemaSrc := genCmd.String("schema_source", "", "[required] Location of the GraphQL schema file. Either absolute or relative.")
	genQuerySrc := genCmd.String("query_source", "", "[required] Location of the GraphQL query file. Either absolute or relative.")
	genPackageName := genCmd.String("package_name", "", "[optional] Name of the generated Go GraphQL client package; defaults to the name of the GraphQL query file with 'grafik_' prefix.")
	genClientName := genCmd.String("client_name", "", "[optional] Name of the generated Go GraphQL client; defaults to the name of the GraphQL query file with 'Grafik' prefix and 'Client' postfix.")
	genDestination := genCmd.String("destination", "./", "[optional] Output filename with path. Either absolute or relative; defaults to the current directory and client name.")
	genUsePointers := genCmd.Bool("use_pointers", false, "[optional] Generate public GraphQL structs' fields as pointers; defaults to false")

	if len(os.Args) < 2 {
		fmt.Println("gen or help subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gen":
		_ = genCmd.Parse(os.Args[2:])
	case "help":
		usage(genCmd)
		os.Exit(0)
	default:
		usage(genCmd)
		os.Exit(1)
	}

	err := validateGenOptions(genSchemaSrc, genQuerySrc)
	if err != nil {
		usage(genCmd)
		os.Exit(1)
	}

	if !genCmd.Parsed() {
		return
	}

	cli := cli{
		schemaSource: genSchemaSrc,
		querySource:  genQuerySrc,
		packageName:  genPackageName,
		clientName:   genClientName,
		destination:  genDestination,
		usePointers:  genUsePointers,
	}

	schemaContent, err := getFileContent(cli.schemaSource)
	if err != nil {
		panic(err)
	}

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(schemaContent),
		Name:  cli.parseSchemaName(),
	})
	print(err.Error())

	queryContent, err := getFileContent(cli.querySource)
	if err != nil {
		panic(err)
	}
	query, err := gqlparser.LoadQuery(schema, string(queryContent))
	print(err.Error())

	additionalInfo := evaluator.AdditionalInfo{
		PackageName: cli.parsePackageName(),
		ClientName:  cli.parseClientName(),
		UsePointers: *genUsePointers,
	}

	e := evaluator.New("./", schema, query, additionalInfo)

	fileName := getFileDestName(additionalInfo.ClientName, cli.destination)
	writeFile(e.Generate(), fileName)
}

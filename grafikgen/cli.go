// package main provides grafikgen CLI tools used for generating grafik clients.
package main

import (
	"flag"
	"github.com/Bartosz-D3V/grafik/evaluator"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"log"
	"os"
)

// cli struct is a wrapper containing all fields required to generate grafik client provided through CLI arguments.
type cli struct {
	schemaSource *string
	querySource  *string
	packageName  *string
	clientName   *string
	destination  *string
	usePointers  *bool
}

func main() {
	genCmd := flag.NewFlagSet("", flag.ExitOnError)
	genSchemaSrc := genCmd.String("schema_source", "", "[required] Location of the GraphQL schema file. Either absolute or relative.")
	genQuerySrc := genCmd.String("query_source", "", "[required] Location of the GraphQL query file. Either absolute or relative.")
	genPackageName := genCmd.String("package_name", "", "[optional] Name of the generated Go GraphQL client package; defaults to the name of the GraphQL query file with 'grafik_' prefix.")
	genClientName := genCmd.String("client_name", "", "[optional] Name of the generated Go GraphQL client; defaults to the name of the GraphQL query file with 'Grafik' prefix and 'Client' postfix.")
	genDestination := genCmd.String("destination", "./", "[optional] Output filename with path. Either absolute or relative; defaults to the current directory and client name.")
	genUsePointers := genCmd.Bool("use_pointers", false, "[optional] Generate public GraphQL structs' fields as pointers; defaults to false.")

	switch os.Args[0] {
	case "help":
		usage(genCmd)
		os.Exit(0)
	default:
		_ = genCmd.Parse(os.Args[1:])
	}

	if !genCmd.Parsed() {
		usage(genCmd)
		os.Exit(1)
	}

	cli := cli{
		schemaSource: genSchemaSrc,
		querySource:  genQuerySrc,
		packageName:  genPackageName,
		clientName:   genClientName,
		destination:  genDestination,
		usePointers:  genUsePointers,
	}

	if *cli.schemaSource == "" || *cli.querySource == "" {
		log.Fatal("grafikgen requires at least two flags - schema_source and query_source.")
	}

	schemaContent, err := getFileContent(cli.schemaSource)
	if err != nil {
		log.Fatalf("Failed to read content of GraphQL schema file. Cause: %s", err.Error())
	}

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(schemaContent),
	})

	// gqlparser returns err that is not nil even when schema is parsed correctly
	if err.Error() != "" {
		log.Fatalf("Failed to parse GraphQL schema file. Cause: %s", err.Error())
	}

	queryContent, err := getFileContent(cli.querySource)
	if err != nil {
		log.Fatalf("Failed to read content of GraphQL query file. Cause: %s", err.Error())
	}

	query, err := gqlparser.LoadQuery(schema, string(queryContent))
	// gqlparser returns err that is not nil even when schema is parsed correctly
	if err.Error() != "" {
		log.Fatalf("Failed to parse GraphQL query file. Cause: %s", err.Error())
	}

	additionalInfo := evaluator.AdditionalInfo{
		PackageName: cli.parsePackageName(),
		ClientName:  cli.parseClientName(),
		UsePointers: *genUsePointers,
	}

	e, err := evaluator.New(schema, query, additionalInfo)
	if err != nil {
		log.Fatal(err)
	}

	fileName := cli.getFileDestName(additionalInfo.ClientName)

	fileContent, err := e.Generate()
	if err != nil {
		log.Fatal(err)
	}
	err = writeFile(fileContent, fileName)
	if err != nil {
		log.Fatal(err)
	}
}

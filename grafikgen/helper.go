package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const usageTxt = `grafik is GraphQL schema based client and generator.

grafikgen is a CLI for generating GraphQL Go clients.

Supported sub-commands:
	help - prints this message

Generate Go GraphQL client by providing location of GraphQL schema and GraphQL queries file.
Example:
	grafikgen -schema_source=./schemas/my_schema.graphql -query_source=./schemas/my_query.graphql [other options]

To customize generated Go GraphQL client provide auxiliary options for client name, package and destination folder.
Example:
	grafikgen -schema_source=./schemas/my_schema.graphql -query_source=./schemas/my_query.graphql -package=app -client_name=MyGraphqlClient -destination=./app/my_client.go

To display this message use help:
Example:
	grafikgen help
`

const rwe = 0755

// writeFile writes the whole slice of bytes to new file.
// Used to create the generated grafik client file.
func writeFile(content []byte, fullDist string) {
	dir, _ := filepath.Split(fullDist)
	if dir != "" {
		err := os.MkdirAll(dir, rwe)
		if err != nil {
			log.Fatalf("Failed to create folder %s. Cause: %s", dir, err.Error())
		}
	}
	openFile, err := os.OpenFile(fullDist, os.O_RDWR|os.O_CREATE|os.O_TRUNC, rwe)
	if err != nil {
		log.Fatalf("Failed to open generated file %s. Cause: %s", fullDist, err.Error())
	}
	_, err = openFile.Write(content)
	if err != nil {
		log.Fatalf("Failed to write content of the grafik client. Cause: %s", err.Error())
	}
	if err := openFile.Close(); err != nil {
		log.Fatalf("Failed to close generated file. Cause: %s", err.Error())
	}
}

// parsePackageName returns name of the package - either defined by user through CLI flag or GraphQL query file name with 'grafik_' prefix.
func (c cli) parsePackageName() string {
	if c.packageName != nil && *c.packageName != "" {
		return *c.packageName
	}
	return fmt.Sprintf("grafik_%s", c.queryFileName())
}

// parseClientName returns name of grafik client - either defined by user through CLI flag or GraphQL schema ame with 'Grafik' prefix and 'Client' suffix.
func (c cli) parseClientName() string {
	if c.clientName != nil && *c.clientName != "" {
		return *c.clientName
	}
	schemaName := c.parseSchemaName()
	return fmt.Sprintf("Grafik%sClient", strings.Title(schemaName))
}

// parseSchemaName returns schema name - either Capitalized client name or Capitalized GraphQL query file source.
func (c cli) parseSchemaName() string {
	if c.clientName != nil && *c.clientName != "" {
		return common.SentenceCase(*c.clientName)
	}
	return c.queryFileName()
}

// queryFileName parses GraphQL query file name.
func (c cli) queryFileName() string {
	baseName := filepath.Base(*c.querySource)
	return common.SentenceCase(strings.Split(baseName, ".")[0])
}

// getFileContent returns content of the file.
func getFileContent(src *string) ([]byte, error) {
	if src == nil || *src == "" {
		return nil, errors.New("provided source is empty")
	}
	if path.IsAbs(*src) {
		return ioutil.ReadFile(*src)
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	absSrc := path.Join(wd, *src)
	return ioutil.ReadFile(absSrc)
}

// getFileDestName returns destination file name - either defined via CLI flag or same as client name.
func getFileDestName(clientName string, dist *string) string {
	if dist == nil || *dist == "" {
		return clientName
	}
	if strings.Contains(*dist, ".go") {
		return *dist
	}
	return fmt.Sprintf("%s.go", filepath.Join(*dist, clientName))
}

// validateGenOptions returns error if any CLI flag is nil or empty.
// Used to validate provided mandatory flags.
func validateGenOptions(opts ...*string) error {
	for _, opt := range opts {
		if opt == nil || *opt == "" {
			return errors.New("insufficient number of command parameters")
		}
	}
	return nil
}

// usage prints help usage text.
func usage(fs *flag.FlagSet) {
	_, _ = io.WriteString(os.Stdout, usageTxt)
	_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Options for %s command: \n", fs.Name()))
	fs.PrintDefaults()
}

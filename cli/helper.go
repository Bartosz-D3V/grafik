package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Bartosz-D3V/grafik/common"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const usageTxt = `grafik is GraphQL schema based client and generator.

Supported commands:
	help - prints this message
	gen  - generates a GraphQL Go client

Generate Go GraphQL client by providing location of GraphQL schema and GraphQL queries file.
Example:
	grafik gen -schema_source=./schemas/my_schema.graphql -query_source=./schemas/my_query.graphql [other options]

To customize generated Go GraphQL client provide auxiliary options for client name, package and destination folder.
Example:
	grafik gen -schema_source=./schemas/my_schema.graphql -query_source=./schemas/my_query.graphql -package=app -client_name=MyGraphqlClient -destination=./app/my_client.go

`

func writeFile(content []byte, fullDist string) {
	dir, _ := filepath.Split(fullDist)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	openFile, err := os.OpenFile(fullDist, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	_, err = openFile.Write(content)
	if err != nil {
		panic(err)
	}
	if err := openFile.Close(); err != nil {
		panic(err)
	}
}

func (c cli) parsePackageName() string {
	if c.packageName != nil && *c.packageName != "" {
		return *c.packageName
	}
	return fmt.Sprintf("grafik_%s", c.queryFileName())
}

func (c cli) parseClientName() string {
	if c.clientName != nil && *c.clientName != "" {
		return *c.clientName
	}
	schemaName := c.parseSchemaName()
	return fmt.Sprintf("Grafik%sClient", strings.Title(schemaName))
}

func (c cli) parseSchemaName() string {
	if c.clientName != nil && *c.clientName != "" {
		return common.SentenceCase(*c.clientName)
	}
	return c.queryFileName()
}

func (c cli) queryFileName() string {
	baseName := filepath.Base(*c.querySource)
	return common.SentenceCase(strings.Split(baseName, ".")[0])
}

func getFileContent(src *string) ([]byte, error) {
	if src == nil || *src == "" {
		return nil, errors.New("provided source is empty")
	}
	if path.IsAbs(*src) {
		return ioutil.ReadFile(*src)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		absSrc := path.Join(wd, *src)
		return ioutil.ReadFile(absSrc)
	}
}

func getFileDestName(clientName string, dist *string) string {
	if dist == nil || *dist == "" {
		return clientName
	}
	dir, file := filepath.Split(*dist)
	if file != "" {
		return *dist
	}
	return fmt.Sprintf("%s.go", filepath.Join(dir, clientName))
}

func validateGenOptions(opts ...*string) error {
	for _, opt := range opts {
		if opt == nil || *opt == "" {
			return errors.New("insufficient number of command parameters")
		}
	}
	return nil
}

func usage(fs *flag.FlagSet) {
	_, _ = io.WriteString(os.Stdout, usageTxt)
	_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Options for %s command: \n", fs.Name()))
	fs.PrintDefaults()
}

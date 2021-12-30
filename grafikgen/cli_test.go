package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCli_parsePackageName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cli cli
		exp string
	}{
		{
			cli{
				packageName: nil,
				querySource: strPtr("apollo_api.graphql"),
			},
			"grafik_apollo_api",
		},
		{
			cli{
				packageName: strPtr(""),
				querySource: strPtr("apollo_api.graphql"),
			},
			"grafik_apollo_api",
		},
		{
			cli{
				packageName: strPtr("apolloApi"),
				querySource: strPtr("apollo_api.graphql"),
			},
			"apolloApi",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.cli.parsePackageName())
	}
}

func TestCli_parseClientName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cli cli
		exp string
	}{
		{
			cli{
				clientName:  nil,
				querySource: strPtr("apollo_api.graphql"),
			},
			"GrafikApolloApiClient",
		},
		{
			cli{
				clientName:  strPtr(""),
				querySource: strPtr("apollo_api.graphql"),
			},
			"GrafikApolloApiClient",
		},
		{
			cli{
				clientName:  strPtr(""),
				querySource: strPtr("apollo-api.graphql"),
			},
			"GrafikApolloApiClient",
		},
		{
			cli{
				clientName:  strPtr("apolloApi"),
				querySource: strPtr("apollo_api.graphql"),
			},
			"apolloApi",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.cli.parseClientName())
	}
}

func TestCli_getFileDestName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cli         cli
		pClientName string
		exp         string
	}{
		{
			cli{destination: strPtr("my_client.go")},
			"MyClient",
			"my_client.go",
		},
		{
			cli{destination: strPtr("/dev/null/client/my_client.go")},
			"MyClient",
			"/dev/null/client/my_client.go",
		},
		{
			cli{destination: strPtr("/dev/null/client/")},
			"MyClient",
			"/dev/null/client/MyClient.go",
		},
		{
			cli{destination: strPtr("/dev/null/client")},
			"MyClient",
			"/dev/null/client/MyClient.go",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, test.cli.getFileDestName(test.pClientName))
	}
}

func strPtr(s string) *string {
	return &s
}

// Generated with grafik. DO NOT EDIT

package main

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type IssueConnection struct {
	TotalCount int `json:"totalCount"`
}

type Language struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type LanguageConnection struct {
	Nodes []Language `json:"nodes"`
}

type Repository struct {
	Forks      RepositoryConnection `json:"forks"`
	Issues     IssueConnection      `json:"issues"`
	Languages  LanguageConnection   `json:"languages"`
	Name       string               `json:"name"`
	Stargazers StargazerConnection  `json:"stargazers"`
	Watchers   UserConnection       `json:"watchers"`
}

type RepositoryConnection struct {
	Edges      []RepositoryEdge `json:"edges"`
	TotalCount int              `json:"totalCount"`
}

type RepositoryEdge struct {
	Node Repository `json:"node"`
}

type StargazerConnection struct {
	TotalCount int `json:"totalCount"`
}

type StarredRepositoryConnection struct {
	TotalCount int `json:"totalCount"`
}

type User struct {
	Login               string                      `json:"login"`
	Repositories        RepositoryConnection        `json:"repositories"`
	StarredRepositories StarredRepositoryConnection `json:"starredRepositories"`
}

type UserConnection struct {
	TotalCount int `json:"totalCount"`
}

const getData = `query getData {
    viewer {
        login
        starredRepositories {
            totalCount
        }
        repositories(first: 3) {
            edges {
                node {
                    name
                    languages(first: 5) {
                        nodes {
                            name,
                            color
                        }
                    }
                    stargazers {
                        totalCount
                    }
                    forks {
                        totalCount
                    }
                    watchers {
                        totalCount
                    }
                    issues(states:[OPEN]) {
                        totalCount
                    }
                }
            }
        }
    }
}`

type GithubClient interface {
	GetData(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *githubClient) GetData(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getData, params, header)
}

type GetDataResponse struct {
	Data   GetDataData    `json:"data"`
	Errors []GraphQLError `json:"errors"`
}

type GetDataData struct {
	Viewer User `json:"viewer"`
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []GraphQLErrorLocation `json:"locations"`
	Extensions GraphQLErrorExtensions `json:"extensions"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type GraphQLErrorExtensions struct {
	Code string `json:"code"`
}

type githubClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) GithubClient {
	return &githubClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}

// client is a package that contains the code used internally by grafik to prepare & send HTTP requests.

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// A Client is an interface that defines a contract of grafik internal GraphQL client.
// It can be mocked with tools like https://github.com/golang/mock in unit tests.
type Client interface {
	Execute(query string, params map[string]interface{}, header *http.Header) (*http.Response, error)
}

// client is a private struct that can be created with New function.
type client struct {
	// endpoint specifies full URL address of the GraphQL endpoint.
	endpoint string

	// httpClient is a pointer to an instance of http.Client. It can be fully customized to provide authentication mechanism, timeout etc.
	httpClient *http.Client
}

// New endpoint creates an instance of the client.
func New(endpoint string, httpClient *http.Client) Client {
	return &client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

// Execute is a receiver function used by generated grafik client to execute HTTP requests.
// Caller method is responsible for closing the body reader.
func (c *client) Execute(query string, params map[string]interface{}, header *http.Header) (*http.Response, error) {
	q := c.formatQuery(query)
	req := GraphQLRequest{
		Query:     q,
		Variables: params,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, GraphQLCallError{"Parsing GraphQL request failed", err.Error()}
	}

	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, GraphQLCallError{"Preparation of GraphQL call failed", err.Error()}
	}

	if header != nil {
		httpReq.Header = *header
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return httpRes, GraphQLCallError{"GraphQL call failed", err.Error()}
	}

	return httpRes, nil
}

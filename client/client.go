package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client interface {
	Execute(query string, params map[string]interface{}, header *http.Header) (*http.Response, error)
}

type client struct {
	endpoint   string
	httpClient *http.Client
}

func New(endpoint string, httpClient *http.Client) Client {
	return &client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

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

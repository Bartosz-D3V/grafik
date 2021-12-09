// Generated with ggrafik. DO NOT EDIT

package main

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type RocketsResult struct {
	Result Result   `json:"result"`
	Data   []Rocket `json:"data"`
}

type Result struct {
	TotalCount int `json:"totalCount"`
}

type Rocket struct {
	TotalPerLaunch int    `json:"total_per_launch"`
	Country        string `json:"country"`
	Name           string `json:"name"`
}

const getRocketResults = `query getRocketResults($limit: Int){
    rocketsResult(limit: $limit) {
        result {
            totalCount
        }
        data {
            name
            country
            total_per_launch
        }
    }
}
`

type UserssClient interface {
	GetRocketResults(limit int, header *http.Header) (*http.Response, error)
}

func (c *userssClient) GetRocketResults(limit int, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["limit"] = limit

	return c.ctrl.Execute(getRocketResults, params, header)
}

type GetRocketResultsResponse struct {
	Data   GetRocketResultsData `json:"data"`
	Errors []Error              `json:"errors"`
}

type GetRocketResultsData struct {
	RocketsResult RocketsResult `json:"rocketsResult"`
}

type Error struct {
	Message    string     `json:"message"`
	Locations  []Location `json:"locations"`
	Extensions Extension  `json:"extensions"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Extension struct {
	Code string `json:"code"`
}

type userssClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) UserssClient {
	return &userssClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
